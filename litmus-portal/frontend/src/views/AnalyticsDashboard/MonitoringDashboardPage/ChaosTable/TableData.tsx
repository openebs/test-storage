import { Typography } from '@material-ui/core';
import React from 'react';
import CheckBox from '../../../../components/CheckBox';
import { ChaosEventDetails } from '../../../../models/dashboardsData';
import useStyles, { StyledTableCell } from './styles';

interface TableDataProps {
  data: ChaosEventDetails;
  itemSelectionStatus: boolean;
  labelIdentifier: string;
}

const TableData: React.FC<TableDataProps> = ({
  data,
  itemSelectionStatus,
  labelIdentifier,
}) => {
  const classes = useStyles();

  return (
    <>
      <StyledTableCell padding="checkbox" className={classes.checkbox}>
        <CheckBox
          checked={itemSelectionStatus}
          inputProps={{ 'aria-labelledby': labelIdentifier }}
        />
      </StyledTableCell>
      <StyledTableCell>
        <div
          className={classes.colorBar}
          style={{
            background: data.legend,
          }}
        />
      </StyledTableCell>
      <StyledTableCell>
        <Typography className={classes.tableObjects}>
          {data.workflow}
        </Typography>
      </StyledTableCell>
      <StyledTableCell>
        <Typography className={classes.tableObjects}>
          {data.experiment}
        </Typography>
      </StyledTableCell>
      <StyledTableCell>
        <Typography className={classes.tableObjects}>{data.target}</Typography>
      </StyledTableCell>
      <StyledTableCell>
        <Typography
          className={`${classes.tableObjects} ${
            data.result === 'Pass'
              ? classes.passColor
              : data.result === 'Fail'
              ? classes.failColor
              : classes.awaitedColor
          }`}
        >
          {data.result}
        </Typography>
      </StyledTableCell>
    </>
  );
};
export default TableData;
