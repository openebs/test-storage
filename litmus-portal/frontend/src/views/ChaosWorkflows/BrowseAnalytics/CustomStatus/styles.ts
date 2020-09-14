import { makeStyles } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  passed: {
    width: '4.75rem',
    textAlign: 'center',
    borderRadius: 3,
    paddingTop: theme.spacing(0.375),
    paddingBottom: theme.spacing(0.375),
    color: theme.palette.primary.dark,
    backgroundColor: 'rgba(16, 155, 103, 0.1)',
    verticalAlign: 'middle',
    display: 'inline-flex',
  },
  failed: {
    width: '4.75rem',
    textAlign: 'center',
    borderRadius: 3,
    paddingTop: theme.spacing(0.375),
    paddingBottom: theme.spacing(0.375),
    color: theme.palette.error.dark,
    backgroundColor: theme.palette.error.light,
    verticalAlign: 'middle',
    display: 'inline-flex',
  },
  statusFont: {
    fontSize: '0.725rem',
  },
  miniIcons: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(0.2),
    display: 'block',
  },
}));

export default useStyles;