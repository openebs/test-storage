import { useQuery } from '@apollo/client';
import { Box, Divider } from '@material-ui/core';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import React, { useCallback, useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { useLocation } from 'react-router-dom';
import CustomBreadCrumbs from '../../../components/BreadCrumbs';
import { Message, NotificationIds } from '../../../models/header';
import {
  CurrentUserDedtailsVars,
  CurrentUserDetails,
  UserData,
} from '../../../models/user';
import { RootState } from '../../../redux/reducers';
import NotificationsDropdown from './NotificationDropdown';
import ProfileDropdownSection from './ProfileDropdownSection';
import useStyles from './styles';
import { GET_USER } from '../../../graphql';
import { Member, Project } from '../../../models/project';
import useActions from '../../../redux/actions';
import * as UserActions from '../../../redux/actions/user';

const Header: React.FC = () => {
  const classes = useStyles();
  const userData: UserData = useSelector((state: RootState) => state.userData);
  const { username } = userData;
  const user = useActions(UserActions);
  // Query to get user details
  const { data } = useQuery<CurrentUserDetails, CurrentUserDedtailsVars>(
    GET_USER,
    { variables: { username } }
  );
  const name: string = data?.getUser.name ?? '';
  const email: string = data?.getUser.email ?? '';
  const projects: Project[] = data?.getUser.projects ?? [];
  const [userRole, setUserRole] = useState<string>(userData.userRole);
  const [selectedProjectName, setSelectedProjectName] = useState<string>(
    userData.selectedProjectName
  );
  const [selectedProject, setSelectedProject] = useState(
    userData.selectedProjectID
  );

  const setSelectedProjectID = (selectedProjectID: string) => {
    const updatedUserDetails = { role: '', projectName: '' };
    setSelectedProject(selectedProjectID);
    projects.forEach((project) => {
      if (selectedProjectID === project.id) {
        const memberList: Member[] = project.members;
        memberList.forEach((member) => {
          if (member.user_name === data?.getUser.username) {
            updatedUserDetails.role = member.role;
            setUserRole(member.role);
          }
        });
        updatedUserDetails.projectName = project.name;
        setSelectedProjectName(project.name);
      }
    });
    user.updateUserDetails({
      selectedProjectID,
      userRole: updatedUserDetails.role,
      selectedProjectName: updatedUserDetails.projectName,
    });
  };

  // Fetch and Set Notifications from backend.

  const [messages, setMessages] = useState<Message[]>([]);

  const [countOfMessages, setCountOfMessages] = useState(0);

  const fetchRandomMessages = useCallback(() => {
    const messages = [];

    const notificationsList = [
      {
        id: '1',
        messageType: 'Pod Delete workflow',
        Message: 'complete',
        generatedTime: '',
      },
      {
        id: '2',
        messageType: 'Argo Chaos workflow',
        Message: 'started started',
        generatedTime: '',
      },
      {
        id: '3',
        messageType: 'New',
        Message: 'crashed',
        generatedTime: '',
      },
    ];

    const iterations = notificationsList.length;

    const oneDaySeconds = 60 * 60 * 24;

    let curUnix = Math.round(
      new Date().getTime() / 1000 - iterations * oneDaySeconds
    );

    for (let i = 0; i < iterations; i += 1) {
      const notificationItem = notificationsList[i];
      const message = {
        sequenceID: (i as unknown) as string,
        id: notificationItem.id,
        messageType: notificationItem.messageType,
        date: curUnix,
        text: `${notificationItem.messageType}- ${notificationItem.Message}`,
      };
      curUnix += oneDaySeconds;
      messages.push(message);
    }
    messages.reverse();
    setMessages(messages);
  }, [setMessages]);

  const deleteNotification = (notificationIDs: NotificationIds) => {
    for (let i = 0; i < messages.length; i += 1) {
      if (messages[i].sequenceID === notificationIDs.sequenceID) {
        if (i > -1) {
          messages.splice(i, 1);
        }
      }
    }
    // send POST request with #notificationIDs.id to update db with notification
    // id marked as disissed from active or persist it in redux or cookie.
    setMessages(messages);
    setCountOfMessages(messages.length);
  };

  useEffect(() => {
    fetchRandomMessages();
  }, [fetchRandomMessages]);

  useEffect(() => {
    setSelectedProject(userData.selectedProjectID);
    setUserRole(userData.userRole);
    setSelectedProjectName(userData.selectedProjectName);
  }, [userData.selectedProjectID]);

  return (
    <div>
      <AppBar position="relative" className={classes.appBar} elevation={0}>
        <Toolbar>
          <div style={{ width: '100%' }}>
            <Box display="flex" p={1} className={classes.headerFlex}>
              <Box p={1} flexGrow={8} className={classes.headerFlexExtraPadded}>
                <CustomBreadCrumbs location={useLocation().pathname} />
              </Box>
              <Box p={1} className={classes.headerFlexPadded}>
                <NotificationsDropdown
                  count={`${countOfMessages}`}
                  messages={messages}
                  CallbackToHeaderOnDeleteNotification={deleteNotification}
                />
              </Box>
              <Box p={1} flexGrow={1} className={classes.headerFlexProfile}>
                <ProfileDropdownSection
                  name={name}
                  email={email}
                  username={username}
                  projects={projects}
                  selectedProjectID={selectedProject}
                  CallbackToSetSelectedProjectID={setSelectedProjectID}
                  selectedProjectName={selectedProjectName}
                  userRole={userRole}
                />
              </Box>
            </Box>
          </div>
        </Toolbar>
        <Divider />
      </AppBar>
    </div>
  );
};

export default Header;
