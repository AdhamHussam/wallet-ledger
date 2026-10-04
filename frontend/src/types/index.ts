export interface Account {
  id: number
  owner: string
  balance: number
  currency: string
  created_at: string
}

export interface Entry {
  id: number
  account_id: number
  amount: number
  created_at: string
}

export interface Transfer {
  id: number
  from_account_id: number
  to_account_id: number
  amount: number
  created_at: string
}

export interface TransferResponse {
  transfer: Transfer
  from_account: Account
  to_account: Account
  from_entry: Entry
  to_entry: Entry
}

export interface CreateAccountRequest {
  owner: string
  currency: string
  initial_balance?: number
}

export interface TransferRequest {
  from_account_id: number
  to_account_id: number
  amount: number
  currency: string
}

export const SUPPORTED_CURRENCIES = ['USD', 'EUR', 'GBP', 'CAD', 'JPY', 'AUD', 'CHF'] as const
