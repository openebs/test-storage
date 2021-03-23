import { makeStyles } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  rootContainer: {
    height: '100%',
    width: '100%',

    background: theme.palette.primary.main,
  },
  rootDiv: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-around',
  },
  rootLitmusText: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
  },
  HeaderText: {
    maxWidth: '23.56rem',
    fontWeight: 500,
    fontSize: '2rem',
    color: theme.palette.text.secondary,
    margin: theme.spacing(1.5, 0, 2.5, 0),
  },
  litmusText: {
    maxWidth: '23.56rem',
    fontWeight: 400,
    fontSize: '1rem',
    color: theme.palette.text.secondary,
  },
  inputDiv: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    maxWidth: '20rem',
    marginTop: theme.spacing(1),
    marginLeft: theme.spacing(7.5),
  },
  inputValue: {
    margin: theme.spacing(1.75, 0, 1.75, 0),
    width: '100%',
  },
  loginButton: {
    marginTop: theme.spacing(1.875),
    background: theme.palette.primary.light,
    color: theme.palette.text.secondary,
    maxWidth: '8rem',
  },
}));
export default useStyles;
