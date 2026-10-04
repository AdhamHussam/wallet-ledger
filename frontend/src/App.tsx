import { useState, useEffect, useCallback } from 'react'
import { api } from './api/client'
import type { Account, Transfer, CreateAccountRequest, TransferRequest } from './types'
import { Navbar } from './components/Navbar'
import { StatsOverview } from './components/StatsOverview'
import { AccountCard } from './components/AccountCard'
import { CreateAccountModal } from './components/CreateAccountModal'
import { TransferModal } from './components/TransferModal'
import { TransferHistory } from './components/TransferHistory'
import { Search, Plus, Filter, AlertCircle, CheckCircle } from 'lucide-react'

export function App() {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [transfers, setTransfers] = useState<Transfer[]>([])
  const [isHealthy, setIsHealthy] = useState<boolean | null>(null)
  const [loading, setLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [currencyFilter, setCurrencyFilter] = useState('ALL')

  // Modals state
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [isTransferOpen, setIsTransferOpen] = useState(false)
  const [transferFromAccount, setTransferFromAccount] = useState<Account | null>(null)
  const [transferToAccount, setTransferToAccount] = useState<Account | null>(null)

  // Notification Toast
  const [toast, setToast] = useState<{ type: 'success' | 'error'; message: string } | null>(null)

  const showToast = (type: 'success' | 'error', message: string) => {
    setToast({ type, message })
    setTimeout(() => setToast(null), 4000)
  }

  const checkHealth = useCallback(async () => {
    try {
      const res = await api.checkHealth()
      setIsHealthy(res.status === 'up')
    } catch {
      setIsHealthy(false)
    }
  }, [])

  const loadAccounts = useCallback(async (isManual = false) => {
    if (isManual) setIsRefreshing(true)
    try {
      const data = await api.listAccounts(1, 100)
      setAccounts(data)
    } catch (err: any) {
      showToast('error', `Failed to load accounts: ${err.message}`)
    } finally {
      setLoading(false)
      setIsRefreshing(false)
    }
  }, [])

  useEffect(() => {
    checkHealth()
    loadAccounts()

    // Poll health every 30 seconds
    const interval = setInterval(checkHealth, 30000)
    return () => clearInterval(interval)
  }, [checkHealth, loadAccounts])

  const handleCreateAccount = async (req: CreateAccountRequest) => {
    const newAccount = await api.createAccount(req)
    setAccounts((prev) => [newAccount, ...prev])
    showToast('success', `Account #${newAccount.id} created for ${newAccount.owner}`)
  }

  const handleTransfer = async (req: TransferRequest) => {
    const res = await api.createTransfer(req)
    // Update local accounts state with returned updated balances
    setAccounts((prev) =>
      prev.map((acc) => {
        if (acc.id === res.from_account.id) return res.from_account
        if (acc.id === res.to_account.id) return res.to_account
        return acc
      })
    )
    // Prepend to transfers audit log
    setTransfers((prev) => [res.transfer, ...prev])
    showToast(
      'success',
      `Transferred ${res.transfer.amount} ${res.from_account.currency} from #${req.from_account_id} to #${req.to_account_id}`
    )
    return res
  }

  const openTransferFrom = (account: Account) => {
    setTransferFromAccount(account)
    setTransferToAccount(null)
    setIsTransferOpen(true)
  }

  const openTransferTo = (account: Account) => {
    setTransferFromAccount(null)
    setTransferToAccount(account)
    setIsTransferOpen(true)
  }

  // Filtered accounts
  const filteredAccounts = accounts.filter((acc) => {
    const matchesSearch =
      acc.owner.toLowerCase().includes(searchQuery.toLowerCase()) ||
      acc.id.toString() === searchQuery.trim()
    const matchesCurrency = currencyFilter === 'ALL' || acc.currency === currencyFilter
    return matchesSearch && matchesCurrency
  })

  const uniqueCurrencies = Array.from(new Set(accounts.map((a) => a.currency)))

  return (
    <div className="min-h-screen bg-[#090d16] text-gray-100 flex flex-col font-sans">
      {/* Toast popup */}
      {toast && (
        <div
          className={`fixed bottom-6 right-6 z-50 px-4 py-3 rounded-xl shadow-2xl flex items-center space-x-3 text-sm font-medium border animate-slide-up ${
            toast.type === 'success'
              ? 'bg-emerald-950/90 border-emerald-500/50 text-emerald-200'
              : 'bg-red-950/90 border-red-500/50 text-red-200'
          }`}
        >
          {toast.type === 'success' ? (
            <CheckCircle className="h-5 w-5 text-emerald-400" />
          ) : (
            <AlertCircle className="h-5 w-5 text-red-400" />
          )}
          <span>{toast.message}</span>
        </div>
      )}

      {/* Navbar */}
      <Navbar
        isHealthy={isHealthy}
        onRefresh={() => loadAccounts(true)}
        onOpenCreateAccount={() => setIsCreateOpen(true)}
        onOpenTransfer={() => {
          setTransferFromAccount(null)
          setTransferToAccount(null)
          setIsTransferOpen(true)
        }}
        isRefreshing={isRefreshing}
      />

      {/* Main Content */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Top metrics bar */}
        <StatsOverview accounts={accounts} transfers={transfers} />

        {/* Accounts Section Header & Filter Toolbar */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <h2 className="text-xl font-bold text-white tracking-tight">Ledger Accounts</h2>
            <p className="text-xs text-gray-400">View balances, audit trails, and manage funds</p>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            {/* Search */}
            <div className="relative">
              <Search className="h-4 w-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="Search owner or ID..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 pr-3.5 py-1.5 bg-gray-900/80 border border-gray-800 rounded-xl text-xs text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition w-44 sm:w-56"
              />
            </div>

            {/* Currency Filter */}
            <div className="flex items-center space-x-1.5 bg-gray-900/80 border border-gray-800 rounded-xl px-2.5 py-1.5 text-xs text-gray-300">
              <Filter className="h-3.5 w-3.5 text-gray-400" />
              <select
                value={currencyFilter}
                onChange={(e) => setCurrencyFilter(e.target.value)}
                className="bg-transparent focus:outline-none cursor-pointer"
              >
                <option value="ALL">All Currencies</option>
                {uniqueCurrencies.map((c) => (
                  <option key={c} value={c} className="bg-gray-900">
                    {c}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>

        {/* Accounts Grid */}
        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-44 rounded-2xl bg-gray-900/40 border border-gray-800 animate-pulse"
              />
            ))}
          </div>
        ) : filteredAccounts.length === 0 ? (
          <div className="bg-gray-900/40 border border-gray-800 rounded-2xl p-12 text-center">
            <p className="text-gray-400 font-medium">No accounts found</p>
            <p className="text-gray-600 text-xs mt-1">
              {searchQuery || currencyFilter !== 'ALL'
                ? 'Try clearing your search or filter'
                : 'Create your first ledger account to start transacting'}
            </p>
            {!searchQuery && currencyFilter === 'ALL' && (
              <button
                onClick={() => setIsCreateOpen(true)}
                className="mt-4 inline-flex items-center space-x-2 px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-xs font-semibold text-white shadow-lg shadow-indigo-600/20 transition"
              >
                <Plus className="h-4 w-4" />
                <span>Create First Account</span>
              </button>
            )}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {filteredAccounts.map((account) => (
              <AccountCard
                key={account.id}
                account={account}
                onTransferFrom={openTransferFrom}
                onTransferTo={openTransferTo}
              />
            ))}
          </div>
        )}

        {/* Audit Log / Transfer History */}
        <div className="mt-12">
          <TransferHistory transfers={transfers} accounts={accounts} />
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800/80 py-6 mt-12 bg-gray-950/40 text-center text-xs text-gray-500">
        Wallet Ledger • ACID Compliant Financial Engine • Go & PostgreSQL
      </footer>

      {/* Modals */}
      <CreateAccountModal
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={handleCreateAccount}
      />

      <TransferModal
        isOpen={isTransferOpen}
        onClose={() => setIsTransferOpen(false)}
        accounts={accounts}
        initialFromAccount={transferFromAccount}
        initialToAccount={transferToAccount}
        onSubmit={handleTransfer}
      />
    </div>
  )
}

export default App
