/* eslint-disable @typescript-eslint/no-unused-vars */
import { Box, IconButton } from '@material-ui/core';
import Avatar from '@material-ui/core/Avatar';
import Typography from '@material-ui/core/Typography';
import ExpandMoreTwoToneIcon from '@material-ui/icons/ExpandMoreTwoTone';
import React, { useRef, useState } from 'react';
import ProfileInfoDropdownItems from './ProfileDropdownItems';
import useStyles from './styles';

interface Project {
  projectName: string;
  statusActive: string;
  id: string;
}

interface CallBackType {
  (selectedProjectID: string): void;
}

interface ProfileInfoDropdownSectionProps {
  name: string;
  email: string;
  username: string;
  projects: Project[];
  selectedProjectID: string;
  CallbackToSetSelectedProjectID: CallBackType;
}

const ProfileDropdownSection = (props: ProfileInfoDropdownSectionProps) => {
  const {
    name,
    email,
    username,
    projects,
    selectedProjectID,
    CallbackToSetSelectedProjectID,
  } = props;

  const classes = useStyles();

  const [isProfilePopoverOpen, setProfilePopoverOpen] = useState(false);

  const profileMenuRef = useRef();

  let initials = ' ';

  if (name) {
    const nameArray = name.split(' ');

    initials =
      nameArray[0][0].toUpperCase() +
      nameArray[nameArray.length - 1][0].toUpperCase();
  }

  const sendSelectedProjectIDToHeader = (selectedProjectID: any) => {
    CallbackToSetSelectedProjectID(selectedProjectID);
    setProfilePopoverOpen(false);
  };

  return (
    <div>
      <Box display="flex" flexDirection="row">
        <Box p={1}>
          {name ? (
            <Avatar
              alt={initials}
              className={classes.avatarBackground}
              style={{ alignContent: 'right' }}
            >
              {initials}
            </Avatar>
          ) : (
            <Avatar
              alt="User"
              className={classes.avatarBackground}
              style={{ alignContent: 'right' }}
            />
          )}
        </Box>

        <Box p={1}>
          <Typography>{name}</Typography>

          <Typography
            style={{
              fontSize: 12,
              color: 'grey',
            }}
          >
            {username}
          </Typography>
        </Box>

        <Box p={1}>
          <IconButton
            edge="end"
            ref={profileMenuRef as any}
            aria-label="account of current user"
            aria-haspopup="true"
            onClick={() => setProfilePopoverOpen(true)}
            style={{ alignContent: 'left' }}
          >
            <ExpandMoreTwoToneIcon htmlColor="grey" />
          </IconButton>
        </Box>
      </Box>
      <ProfileInfoDropdownItems
        anchorEl={profileMenuRef.current as any}
        isOpen={isProfilePopoverOpen}
        onClose={() => setProfilePopoverOpen(false)}
        name={name}
        email={email}
        projects={projects}
        selectedProjectID={selectedProjectID}
        CallbackToSetSelectedProjectIDOnProfileDropdown={
          sendSelectedProjectIDToHeader
        }
      />
    </div>
  );
};

export default ProfileDropdownSection;
