import React, { useState, useEffect } from 'react'
import type { Account, TransferRequest, TransferResponse } from '../types'
import { X, ArrowRightLeft, AlertCircle, CheckCircle2, ArrowRight } from 'lucide-react'

interface TransferModalProps {
  isOpen: boolean
  onClose: () => void
  accounts: Account[]
  initialFromAccount?: Account | null
  initialToAccount?: Account | null
  onSubmit: (req: TransferRequest) => Promise<TransferResponse>
}

export const TransferModal: React.FC<TransferModalProps> = ({
  isOpen,
  onClose,
  accounts,
  initialFromAccount,
  initialToAccount,
  onSubmit,
}) => {
  const [fromAccountId, setFromAccountId] = useState<number | ''>('')
  const [toAccountId, setToAccountId] = useState<number | ''>('')
  const [amount, setAmount] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [successResult, setSuccessResult] = useState<TransferResponse | null>(null)

  useEffect(() => {
    if (isOpen) {
      setError(null)
      setSuccessResult(null)
      setAmount('')
      if (initialFromAccount) {
        setFromAccountId(initialFromAccount.id)
      } else if (accounts.length > 0) {
        setFromAccountId(accounts[0].id)
      }

      if (initialToAccount) {
        setToAccountId(initialToAccount.id)
      } else if (accounts.length > 1) {
        // pick second account if different
        setToAccountId(accounts[1].id)
      }
    }
  }, [isOpen, initialFromAccount, initialToAccount, accounts])

  if (!isOpen) return null

  const fromAccount = accounts.find((a) => a.id === fromAccountId)
  const eligibleReceivers = accounts.filter(
    (a) => a.id !== fromAccountId && (!fromAccount || a.currency === fromAccount.currency)
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!fromAccountId || !toAccountId) {
      setError('Please select both sender and receiver accounts')
      return
    }

    if (fromAccountId === toAccountId) {
      setError('Cannot transfer money to the same account')
      return
    }

    const transferAmount = parseInt(amount, 10)
    if (isNaN(transferAmount) || transferAmount <= 0) {
      setError('Transfer amount must be greater than zero')
      return
    }

    if (fromAccount && fromAccount.balance < transferAmount) {
      setError(`Insufficient funds! Available balance is ${fromAccount.balance} ${fromAccount.currency}`)
      return
    }

    try {
      setLoading(true)
      const result = await onSubmit({
        from_account_id: Number(fromAccountId),
        to_account_id: Number(toAccountId),
        amount: transferAmount,
        currency: fromAccount?.currency || 'USD',
      })
      setSuccessResult(result)
    } catch (err: any) {
      setError(err.message || 'Transfer failed')
    } finally {
      setLoading(false)
    }
  }

  const handleSetMax = () => {
    if (fromAccount && fromAccount.balance > 0) {
      setAmount(fromAccount.balance.toString())
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fade-in">
      <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl relative">
        {/* Header */}
        <div className="flex items-center justify-between pb-4 border-b border-gray-800">
          <div className="flex items-center space-x-2">
            <div className="p-2 bg-indigo-500/10 rounded-lg text-indigo-400">
              <ArrowRightLeft className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white">Execute Money Transfer</h3>
              <p className="text-xs text-gray-400">Atomic, deadlock-safe ACID transaction</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white p-1 rounded-lg hover:bg-gray-800 transition"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Success View */}
        {successResult ? (
          <div className="py-6 text-center space-y-4">
            <div className="mx-auto w-12 h-12 bg-emerald-500/10 text-emerald-400 rounded-full flex items-center justify-center border border-emerald-500/20">
              <CheckCircle2 className="h-7 w-7" />
            </div>

            <div>
              <h4 className="text-xl font-bold text-white">Transfer Completed</h4>
              <p className="text-sm text-gray-400 mt-1">
                Transaction #{successResult.transfer.id} committed to ledger
              </p>
            </div>

            {/* Transfer breakdown */}
            <div className="bg-gray-950/80 border border-gray-800 rounded-xl p-4 text-left space-y-3 font-mono text-xs">
              <div className="flex justify-between items-center text-gray-400">
                <span>Transferred:</span>
                <span className="text-emerald-400 font-bold text-sm">
                  {successResult.transfer.amount} {successResult.from_account.currency}
                </span>
              </div>
              <div className="flex justify-between items-center text-gray-400">
                <span>Sender (Debit):</span>
                <span className="text-red-400">
                  {successResult.from_entry.amount} (Balance: {successResult.from_account.balance})
                </span>
              </div>
              <div className="flex justify-between items-center text-gray-400">
                <span>Receiver (Credit):</span>
                <span className="text-emerald-400">
                  +{successResult.to_entry.amount} (Balance: {successResult.to_account.balance})
                </span>
              </div>
            </div>

            <button
              onClick={onClose}
              className="w-full py-2.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 font-semibold text-white text-sm transition"
            >
              Done
            </button>
          </div>
        ) : (
          /* Form View */
          <form onSubmit={handleSubmit} className="mt-5 space-y-4">
            {error && (
              <div className="p-3 bg-red-500/10 border border-red-500/30 rounded-xl flex items-start space-x-2 text-red-400 text-sm">
                <AlertCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
                <span>{error}</span>
              </div>
            )}

            {/* From Account */}
            <div>
              <div className="flex justify-between items-center mb-1.5">
                <label className="text-xs font-semibold text-gray-300 uppercase tracking-wider">
                  From Account (Debit)
                </label>
                {fromAccount && (
                  <span className="text-xs text-indigo-400 font-medium">
                    Available: {fromAccount.balance} {fromAccount.currency}
                  </span>
                )}
              </div>
              <select
                value={fromAccountId}
                onChange={(e) => setFromAccountId(Number(e.target.value))}
                className="w-full px-3.5 py-2.5 bg-gray-950 border border-gray-800 rounded-xl text-white focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition text-sm"
              >
                {accounts.map((acc) => (
                  <option key={acc.id} value={acc.id}>
                    #{acc.id} — {acc.owner} ({acc.balance} {acc.currency})
                  </option>
                ))}
              </select>
            </div>

            {/* To Account */}
            <div>
              <label className="block text-xs font-semibold text-gray-300 uppercase tracking-wider mb-1.5">
                To Account (Credit)
              </label>
              <select
                value={toAccountId}
                onChange={(e) => setToAccountId(Number(e.target.value))}
                className="w-full px-3.5 py-2.5 bg-gray-950 border border-gray-800 rounded-xl text-white focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition text-sm"
              >
                <option value="" disabled>
                  Select destination account
                </option>
                {eligibleReceivers.map((acc) => (
                  <option key={acc.id} value={acc.id}>
                    #{acc.id} — {acc.owner} ({acc.balance} {acc.currency})
                  </option>
                ))}
              </select>
              {eligibleReceivers.length === 0 && fromAccount && (
                <p className="text-xs text-amber-400 mt-1">
                  No other accounts found with matching currency ({fromAccount.currency}).
                </p>
              )}
            </div>

            {/* Amount */}
            <div>
              <div className="flex justify-between items-center mb-1.5">
                <label className="text-xs font-semibold text-gray-300 uppercase tracking-wider">
                  Transfer Amount
                </label>
                {fromAccount && (
                  <button
                    type="button"
                    onClick={handleSetMax}
                    className="text-xs text-indigo-400 hover:text-indigo-300 font-semibold"
                  >
                    Send Max ({fromAccount.balance})
                  </button>
                )}
              </div>
              <div className="relative">
                <input
                  type="number"
                  min="1"
                  step="1"
                  placeholder="0"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  className="w-full px-3.5 py-2.5 bg-gray-950 border border-gray-800 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition text-sm font-mono pr-16"
                />
                <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-xs font-semibold text-gray-500">
                  {fromAccount?.currency || 'USD'}
                </div>
              </div>
            </div>

            {/* Visual Transfer Pathway */}
            {fromAccountId && toAccountId && fromAccountId !== toAccountId && (
              <div className="p-3 bg-gray-950/60 rounded-xl border border-gray-800/80 flex items-center justify-between text-xs text-gray-400">
                <span className="font-semibold text-white">
                  #{fromAccountId} ({fromAccount?.owner})
                </span>
                <div className="flex items-center space-x-1 text-indigo-400">
                  <span className="font-mono">{amount || '0'}</span>
                  <ArrowRight className="h-4 w-4" />
                </div>
                <span className="font-semibold text-white">
                  #{toAccountId} ({accounts.find((a) => a.id === toAccountId)?.owner})
                </span>
              </div>
            )}

            {/* Actions */}
            <div className="pt-4 flex items-center justify-end space-x-3 border-t border-gray-800 mt-6">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 rounded-xl text-sm font-medium text-gray-300 hover:text-white hover:bg-gray-800 transition"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading || eligibleReceivers.length === 0}
                className="px-5 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-sm font-semibold text-white shadow-lg shadow-indigo-600/30 transition disabled:opacity-50"
              >
                {loading ? 'Processing...' : 'Transfer Funds'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  )
}
