const DATE_TIME_PAD_LENGTH = 2;
const MONTH_INDEX_OFFSET = 1;

export function formatDateTime(value: string): string {
  if (!value) {
    return '';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  const year = date.getFullYear();
  const month = date.getMonth() + MONTH_INDEX_OFFSET;
  const day = date.getDate();
  const hour = String(date.getHours()).padStart(DATE_TIME_PAD_LENGTH, '0');
  const minute = String(date.getMinutes()).padStart(DATE_TIME_PAD_LENGTH, '0');

  return `${String(year)}/${String(month)}/${String(day)} ${hour}:${minute}`;
}
