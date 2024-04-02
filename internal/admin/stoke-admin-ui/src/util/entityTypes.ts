export interface User {
  id: number,
  fname: string,
  lname: string,
  username: string,
  email: string,
  created_at: string,
  provider: string,
}

export interface UserWithCreds {
  id: number,
  fname: string,
  lname: string,
  username: string,
  email: string,
  created_at: string,
  password: string,
  provider: string,
}

export interface Group {
  id: number,
  name: string,
  description: string,
  is_user_group: boolean,
}

export interface Claim {
  id: number
  name: string
  short_name: string
  value: string
  description: string
}
