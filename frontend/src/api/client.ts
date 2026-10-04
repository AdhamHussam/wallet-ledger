import type { Account, CreateAccountRequest, Transfer, TransferRequest, TransferResponse } from '../types'

const API_BASE = ''

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    let errorMsg = `HTTP Error ${res.status}`
    try {
      const data = await res.json()
      if (data && data.error) {
        errorMsg = data.error
      }
    } catch {
      // ignore json parse error
    }
    throw new Error(errorMsg)
  }
  return res.json()
}

export const api = {
  async checkHealth(): Promise<{ status: string; service: string }> {
    const res = await fetch(`${API_BASE}/health`)
    return handleResponse(res)
  },

  async listAccounts(pageId = 1, pageSize = 100): Promise<Account[]> {
    const res = await fetch(`${API_BASE}/accounts?page_id=${pageId}&page_size=${pageSize}`)
    return handleResponse(res)
  },

  async getAccount(id: number): Promise<Account> {
    const res = await fetch(`${API_BASE}/accounts/${id}`)
    return handleResponse(res)
  },

  async createAccount(req: CreateAccountRequest): Promise<Account> {
    const res = await fetch(`${API_BASE}/accounts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req),
    })
    return handleResponse(res)
  },

  async listTransfers(pageId = 1, pageSize = 50): Promise<Transfer[]> {
    const res = await fetch(`${API_BASE}/transfers?page_id=${pageId}&page_size=${pageSize}`)
    return handleResponse(res)
  },

  async createTransfer(req: TransferRequest): Promise<TransferResponse> {
    const res = await fetch(`${API_BASE}/transfers`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req),
    })
    return handleResponse(res)
  },

  async getTransfer(id: number): Promise<Transfer> {
    const res = await fetch(`${API_BASE}/transfers/${id}`)
    return handleResponse(res)
  },
}
