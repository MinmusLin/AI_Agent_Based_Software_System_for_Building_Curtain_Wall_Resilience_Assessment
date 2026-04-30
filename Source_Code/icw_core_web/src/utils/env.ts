export function getRequiredEnv(key: string): string {
  const value = import.meta.env[key] as unknown;
  if (typeof value !== 'string' || !value) {
    throw new Error(`${key} is required`);
  }
  return value;
}
