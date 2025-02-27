import React, { useEffect } from "react";
import { Navbar, Nav } from "react-bootstrap";
import { NavLink, useHistory } from "react-router-dom";

import SearchField, { SearchType } from "src/components/searchField";
import { getPlatformURL, getCredentialsSetting } from "src/utils/createClient";
import { isAdmin, userHref, setCachedUser } from "src/utils";
import { useAuth } from "src/hooks";
import {
  ROUTE_SCENES,
  ROUTE_PERFORMERS,
  ROUTE_TAGS,
  ROUTE_STUDIOS,
  ROUTE_EDITS,
  ROUTE_LOGOUT,
  ROUTE_LOGIN,
  ROUTE_USERS,
  ROUTE_ACTIVATE,
  ROUTE_RESET_PASSWORD,
  ROUTE_HOME,
  ROUTE_REGISTER,
  ROUTE_FORGOT_PASSWORD,
} from "src/constants/route";
import AuthContext from "./AuthContext";

const Main: React.FC = ({ children }) => {
  const history = useHistory();
  const { loading, user } = useAuth();

  useEffect(() => {
    if (loading || user) return;

    if (
      history.location.pathname !== ROUTE_ACTIVATE &&
      history.location.pathname !== ROUTE_REGISTER &&
      history.location.pathname !== ROUTE_LOGIN &&
      history.location.pathname !== ROUTE_FORGOT_PASSWORD &&
      history.location.pathname !== ROUTE_RESET_PASSWORD
    ) {
      history.push(ROUTE_LOGIN);
    }
  }, [loading, user, history]);

  const contextValue = {
    authenticated: user !== undefined,
    user,
  };

  if (!contextValue.authenticated)
    return (
      <AuthContext.Provider value={contextValue}>
        {children}
      </AuthContext.Provider>
    );

  const handleLogout = async () => {
    const res = await fetch(`${getPlatformURL()}logout`, {
      credentials: getCredentialsSetting(),
    });
    setCachedUser();
    if (res.ok) window.location.href = ROUTE_LOGIN;
    return false;
  };

  const renderUserNav = () =>
    contextValue.authenticated &&
    contextValue.user && (
      <>
        <span>Logged in as</span>
        <NavLink
          to={userHref(contextValue.user)}
          className="nav-link ml-auto mr-2"
        >
          {contextValue.user.name}
        </NavLink>
        {isAdmin(user) && (
          <NavLink exact to={ROUTE_USERS} className="nav-link">
            Users
          </NavLink>
        )}
        <NavLink
          to={ROUTE_LOGOUT}
          onClick={handleLogout}
          className="nav-link mr-4"
        >
          Logout
        </NavLink>
      </>
    );

  return (
    <div>
      <Navbar bg="dark" variant="dark" className="px-4">
        <Nav className="row mr-auto">
          <NavLink exact to={ROUTE_HOME} className="nav-link">
            Home
          </NavLink>
          <NavLink to={ROUTE_PERFORMERS} className="nav-link">
            Performers
          </NavLink>
          <NavLink to={ROUTE_SCENES} className="nav-link">
            Scenes
          </NavLink>
          <NavLink to={ROUTE_STUDIOS} className="nav-link">
            Studios
          </NavLink>
          <NavLink to={ROUTE_TAGS} className="nav-link">
            Tags
          </NavLink>
          <NavLink to={ROUTE_EDITS} className="nav-link">
            Edits
          </NavLink>
        </Nav>
        <Nav className="align-items-center">
          {contextValue.authenticated && renderUserNav()}
          <SearchField searchType={SearchType.Combined} navigate showAllLink />
        </Nav>
      </Navbar>
      <div className="StashDBContent container-fluid">
        <AuthContext.Provider value={contextValue}>
          {children}
        </AuthContext.Provider>
      </div>
    </div>
  );
};

export default Main;
