import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

import { formatTimestamp } from './utils';

const styles = {
  table: {
    width: '100%',
    overflowX: 'auto',
  },
  row: {
    height: '0',
  },
  cell: {
    fontSize: '10pt',
    borderBottomStyle: 'none',
    paddingLeft:'0px',
  },
  key: {
    width: '20%',
  },
  value: {
    width: '80%',
  },
};

class NodeInfoRaw extends React.Component {

  renderRow(key, val, formatFunc) {
    if (typeof formatFunc === "function") {
      val = formatFunc(val)
    } else if (formatFunc) {
      console.log("renderRow: formatFunc was not a function:", typeof(formatFunc))
    }
    const { classes } = this.props;
    if ( key && val ) {
      return (
        <TableRow key={key} className={classes.row}>
          <TableCell className={[classes.cell, classes.key]}><b>{key}</b></TableCell>
          <TableCell className={[classes.cell, classes.value]}>{val}</TableCell>
        </TableRow>
      )
    }
    return
  }

  resourceString(resources) {
    if ( resources ) {
      const r = resources
      var s = r.cpuCores + " CPU cores";
      if (r.ramGb) {
        s += ", " + r.ramGb + " GB RAM";
      }
      if (r.diskGb) {
        s += ", " + r.diskGb + " GB disk space";
      }
      if (r.preemptible) {
        s += ", preemptible";
      }
      return s
    }
    return
  }

  renderNode(node) {
    const { classes } = this.props;
    if (!node) {
      return
    }
    return(
      <div>
        <Table className={classes.table}>
          <TableBody>
            {this.renderRow('ID', node.id)}
            {this.renderRow('Hostname', node.hostname)}
            {this.renderRow('State', node.state)}
            {this.renderRow('Resources', node.resources, this.resourceString)}
            {this.renderRow('Available', node.available, this.resourceString)}
            {this.renderRow('Last Ping', node.lastPing, formatTimestamp)}
            {this.renderRow('Version', node.version)}
            {this.renderRow('Tasks', node.task_ids.map(tid => ( 
              <a href={"/tasks/" + tid}>{tid}</a>
            )))}
          </TableBody>
        </Table>
      </div>
    )
  }

  render() {
    const node = this.props.node
    console.log("node", node)
    return (
      <div>
        {this.renderNode(node)}
      </div>
    )
  }
}

NodeInfoRaw.propTypes = {
  node: PropTypes.object.isRequired,
};

const NodeInfo = withStyles(styles)(NodeInfoRaw);
export { NodeInfo };
