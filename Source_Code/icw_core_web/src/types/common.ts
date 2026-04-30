export interface ApiEnvelope<T> {
  code: string;
  message: string;
  data: T;
}

export interface User {
  id: number;
  email: string;
  name: string;
}
