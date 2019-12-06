import React from 'react';
import { Route, Redirect, Switch, withRouter } from "react-router-dom";
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Drawer from '@material-ui/core/Drawer';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import List from '@material-ui/core/List';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import { NavListItems, TaskFilters } from './listItems';
import { TaskList, Task, Node, NodeList, ServiceInfo, NoMatch } from './Pages.js';

const drawerWidth = 240;

const styles = theme => ({
  root: {
    display: 'flex',
  },
  toolbar: {
    paddingRight: 24, // keep right padding when drawer closed
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
  appBar: {
    backgroundColor: "#000000",
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginLeft: 12,
    marginRight: 36,
  },
  menuButtonHidden: {
    display: 'none',
  },
  title: {
    flexGrow: 1,
  },
  drawerPaper: {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: drawerWidth,
    overflowx: "wrap",
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerPaperClose: {
    overflowX: 'hidden',
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    width: 0,
  },
  appBarSpacer: theme.mixins.toolbar,
  content: {
    flexGrow: 1,
    padding: theme.spacing.unit * 3,
    height: '100vh',
    overflow: 'auto',
  },
  chartContainer: {
    marginLeft: -22,
  },
  h5: {
    marginBottom: theme.spacing.unit * 2,
  },
});

class Dashboard extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      open: true,
      stateFilter: "",
      tagsFilter: [{key:"", value: ""}],
    };
  };

  handleDrawerOpen = () => {
    this.setState({ open: true });
  };

  handleDrawerClose = () => {
    this.setState({ open: false });
  };

  updateStateFilter = (event) => {
    this.setState({ stateFilter: event.target.value });
  };

  updateTagsFilter = (index, target, val) => {
    var tags = JSON.parse(JSON.stringify(this.state.tagsFilter));
    if (target === "key") {
      tags[index].key = val;
    } else if (target === "value") {
      tags[index].value = val;
    };
    this.setState({tagsFilter: tags});
  };

  addTagFilter = () => {
    this.setState(state => { 
      const tagsFilter = state.tagsFilter.concat({key: "", value: ""});
      return {tagsFilter};
    });
  };

  removeTagFilter = (index) => {
    this.setState(state => { 
      const tagsFilter = state.tagsFilter.filter((item, j) => index !== j);
      return {tagsFilter};
    });
  };

  render() {
    const { classes } = this.props;
    //console.log("Dashboard props:", this.props)
    //console.log("Dashboard state:", this.state)
    //console.log("Dashboard stateFilter:", this.state.stateFilter)
    //console.log("Dashboard tagsFilter:", this.state.tagsFilter[0])
    return (
      <div className={classes.root}>
        <CssBaseline />
        <AppBar
          position="fixed"
          className={classNames(classes.appBar, this.state.open && classes.appBarShift)}
        >
          <Toolbar disableGutters={!this.state.open} className={classes.toolbar}>
            <IconButton
              color="inherit"
              aria-label="Open drawer"
              onClick={this.handleDrawerOpen}
              className={classNames(
                classes.menuButton,
                this.state.open && classes.menuButtonHidden,
              )}
            >
              <MenuIcon />
            </IconButton>
            <Typography
              component="h1"
              variant="h6"
              color="inherit"
              noWrap
              className={classes.title}
            >
              Funnel
            </Typography>
          </Toolbar>
        </AppBar>
        <Drawer
          variant="permanent"
          classes={{
            paper: classNames(classes.drawerPaper, !this.state.open && classes.drawerPaperClose),
          }}
          open={this.state.open}
        >
          <div className={classes.toolbarIcon}>
            <IconButton onClick={this.handleDrawerClose}>
              <ChevronLeftIcon />
            </IconButton>
          </div>
          <Divider />
          <List>{NavListItems}</List>
          <Divider />
          <TaskFilters 
            show={window.location.pathname.endsWith("tasks") || window.location.pathname === "/"}
            stateFilter={this.state.stateFilter}
            tagsFilter={this.state.tagsFilter}
            updateState={this.updateStateFilter}
            updateTags={this.updateTagsFilter}
            addTag={this.addTagFilter}
            removeTag={this.removeTagFilter}
          />
        </Drawer>
        <main className={classes.content}>
        <div className={classes.appBarSpacer} />
          <Switch>
            <Redirect exact from="/" to="/tasks" />
            <Route exact path="/v1/tasks" render={ (props) => <TaskList {...props} stateFilter={this.state.stateFilter} tagsFilter={this.state.tagsFilter} /> } />
            <Route exact path="/tasks" render={ (props) => <TaskList {...props} stateFilter={this.state.stateFilter} tagsFilter={this.state.tagsFilter} /> } />
            <Route exact path="/v1/tasks/:task_id" component={Task} />
            <Route exact path="/tasks/:task_id" component={Task} />
            <Route exact path="/v1/nodes" render={ (props) => <NodeList {...props} /> } />
            <Route exact path="/nodes" render={ (props) => <NodeList {...props} /> } />
            <Route exact path="/v1/nodes/:node_id" component={Node} />
            <Route exact path="/nodes/:node_id" component={Node} />
            <Route exact path="/v1/tasks/service-info" component={ServiceInfo} />
            <Route exact path="/tasks/service-info" component={ServiceInfo} />
            <Route exact path="/service-info" component={ServiceInfo} />
            <Route component={NoMatch} />
          </Switch>
        </main>
      </div>
    );
  }
}

Dashboard.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withRouter(withStyles(styles)(Dashboard));