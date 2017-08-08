package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"os"
	"path/filepath"
	"strings"
)

// S3Protocol defines the expected URL prefix for S3, "s3://"
const S3Protocol = "s3://"

// S3Backend provides access to an S3 object store.
type S3Backend struct {
	sess *session.Session
}

// NewS3Backend creates an S3Backend session instance
func NewS3Backend(conf config.S3Storage) (*S3Backend, error) {

	// Initialize a session object.
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &S3Backend{sess}, nil
}

// Get copies an object from S3 to the host path.
func (s3b *S3Backend) Get(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting download", "url", url)

	path := strings.TrimPrefix(url, S3Protocol)
	split := strings.SplitN(path, "/", 2)
	bucket := split[0]
	key := split[1]

	var err error

	region, err := s3manager.GetBucketRegion(context.Background(), s3b.sess, bucket, "us-east-1")
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return fmt.Errorf("unable to find bucket %s's region not found\n", bucket)
		} else {
			return err
		}
	}

	// Create a downloader with the session and default options
	sess := s3b.sess.Copy(&aws.Config{Region: aws.String(region)})
	client := s3.New(sess)
	manager := s3manager.NewDownloader(sess)

	switch class {
	case File:
		// Create a file to write the S3 Object contents to.
		hostFile, oerr := os.Create(hostPath)
		if oerr != nil {
			return fmt.Errorf("failed to create file %q, %v", hostPath, err)
		}
		defer hostFile.Close()

		_, err = manager.Download(hostFile, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	case Directory:
		d := directoryDownloader{bucket: bucket, dir: hostPath, Downloader: manager}
		err = client.ListObjectsPages(
			&s3.ListObjectsInput{Bucket: &bucket, Prefix: &key},
			d.eachPage,
		)

	default:
		err = fmt.Errorf("Unknown file class: %s", class)
	}

	if err != nil {
		return err
	}
	log.Info("Finished download", "url", url, "hostPath", hostPath)
	return nil
}

// Put copies an object (file) from the host path to S3.
func (s3b *S3Backend) Put(ctx context.Context, url string, hostPath string, class tes.FileType) ([]*tes.OutputFileLog, error) {
	log.Info("Starting upload", "url", url, "hostPath", hostPath)

	path := strings.TrimPrefix(url, S3Protocol)
	split := strings.SplitN(path, "/", 2)
	bucket := split[0]
	key := split[1]

	region, err := s3manager.GetBucketRegion(context.Background(), s3b.sess, bucket, "us-east-1")
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return nil, fmt.Errorf("unable to find bucket %s's region not found\n", bucket)
		} else {
			return nil, err
		}
	}

	// Create a uploader with the session and default options
	sess := s3b.sess.Copy(&aws.Config{Region: aws.String(region)})
	manager := s3manager.NewUploader(sess)

	var out []*tes.OutputFileLog

	switch class {
	case File:
		f, err := os.Open(hostPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %q, %v", hostPath, err)
		}
		defer f.Close()
		_, err = manager.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   f,
		})
		if err != nil {
			return nil, err
		}

		out = append(out, &tes.OutputFileLog{
			Url:       url,
			Path:      hostPath,
			SizeBytes: fileSize(hostPath),
		})

	case Directory:
		files, err := walkFiles(hostPath)
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			u := url + "/" + f.rel
			fh, err := os.Open(f.abs)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %q, %v", f.abs, err)
			}
			defer fh.Close()
			_, err = manager.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   fh,
			})
			if err != nil {
				return nil, err
			}
			out = append(out, &tes.OutputFileLog{
				Url:       u,
				Path:      f.abs,
				SizeBytes: f.size,
			})
		}

	default:
		return nil, fmt.Errorf("Unknown file class: %s", class)
	}

	log.Info("Finished upload", "url", url, "hostPath", hostPath)
	return out, nil
}

// Supports indicates whether this backend supports the given storage request.
// For S3, the url must start with "s3://".
func (s3b *S3Backend) Supports(url string, hostPath string, class tes.FileType) bool {
	return strings.HasPrefix(url, S3Protocol)
}

type directoryDownloader struct {
	*s3manager.Downloader
	bucket string
	dir    string
}

func (d *directoryDownloader) eachPage(page *s3.ListObjectsOutput, more bool) bool {
	for _, obj := range page.Contents {
		err := d.downloadToFile(*obj.Key)
		if err != nil {
			panic(err)
		}
	}
	return true
}

func (d *directoryDownloader) downloadToFile(key string) error {
	// Create the directories in the path
	file := filepath.Join(d.dir, key)
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return err
	}

	// Setup the local file
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Download the file using the AWS SDK
	_, err = d.Download(fd, &s3.GetObjectInput{
		Bucket: &d.bucket,
		Key:    &key,
	})
	return err
}
