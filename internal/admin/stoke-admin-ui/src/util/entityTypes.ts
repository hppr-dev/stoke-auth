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
  id: number
  fname: string
  lname: string
  username: string
  email: string
  created_at: string
  password: string
  provider: string
}

export interface Group {
  id: number
  name: string
  description: string
}

export interface Claim {
  id: number
  name: string
  short_name: string
  value: string
  description: string
}

export interface GroupLink {
  id: number
  type: string
  resource_spec: string
  claim_group: number
}
