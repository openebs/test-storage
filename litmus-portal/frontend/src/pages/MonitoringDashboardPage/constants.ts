export const PROMETHEUS_ERROR_QUERY_RESOLUTION_LIMIT_REACHED: string = `bad_data: exceeded maximum resolution of 11,000 points per timeseries. Try decreasing the query resolution (?step=XX)`;
export const DEFAULT_CHAOS_EVENT_PROMETHEUS_QUERY_RESOLUTION: string = '1/2';
export const CHAOS_EXPERIMENT_VERDICT_PASS: string = 'Pass';
export const CHAOS_EXPERIMENT_VERDICT_FAIL: string = 'Fail';
export const STATUS_RUNNING: string = 'Running';
export const PROMETHEUS_QUERY_RESOLUTION_LIMIT: number = 11000;
export const MAX_REFRESH_RATE: number = 2147483647;
export const DEFAULT_REFRESH_RATE: number = 10000;
export const ACTIVE: string = 'Active';
export const MINIMUM_TOLERANCE_LIMIT: number = 4;
export const DEFAULT_TOLERANCE_LIMIT: number = 14;
export const INVALID_RESILIENCE_SCORE_STRING: string = 'NaN';
export const DEFAULT_METRIC_SERIES_NAME: string = 'metric';
export const DEFAULT_CHAOS_EVENT_NAME: string = 'chaos';
export const DEFAULT_RELATIVE_TIME_RANGE: number = 1800;
export const DASHBOARD_TYPE_1: string = 'Kubernetes Platform';
export const DASHBOARD_TYPE_2: string = 'Sock Shop';
