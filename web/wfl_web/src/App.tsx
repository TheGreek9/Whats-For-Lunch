import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";

import { GroupView } from "./components/GroupView"
import { CreateGroupController } from "./components/CreateGroupController"


export interface RouteProps {
    history: object
    location: object
    match: {path: string, url: string, isExact: boolean, params: {groupId: string}}
}

export default function App() {
  return (
   <Router>
     <Switch>
      <Route exact strict path="/group/create" render={
         props => <CreateGroupController {...props} />
      }/>
      <Route exact strict path="/group/:groupId" render={
        props => <GroupView {...props}/>
      }/>
      <Route exact strict path="/groups/:groupId" render={
        props => <GroupView {...props}/>
      }/>
     </Switch>
   </Router>
  );
}
