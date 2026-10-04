import React from 'react'
import type { Account } from '../types'
import { ArrowUpRight, ArrowDownLeft, Calendar } from 'lucide-react'

interface AccountCardProps {
  account: Account
  onTransferFrom: (account: Account) => void
  onTransferTo: (account: Account) => void
}

const currencySymbols: Record<string, string> = {
  USD: '$',
  EUR: '€',
  GBP: '£',
  CAD: 'CA$',
  JPY: '¥',
  AUD: 'AU$',
  CHF: 'CHF ',
}

export const AccountCard: React.FC<AccountCardProps> = ({
  account,
  onTransferFrom,
  onTransferTo,
}) => {
  const symbol = currencySymbols[account.currency] || `${account.currency} `
  const dateFormatted = new Date(account.created_at).toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })

  return (
    <div className="bg-gradient-to-b from-gray-900/90 to-gray-950/90 border border-gray-800 hover:border-gray-700/80 rounded-2xl p-5 shadow-lg shadow-black/20 transition group">
      {/* Top Header */}
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center space-x-2">
            <span className="font-semibold text-white text-base tracking-tight">{account.owner}</span>
            <span className="text-xs px-2 py-0.5 rounded-full bg-gray-800 text-gray-400 font-mono">
              #{account.id}
            </span>
          </div>
          <div className="flex items-center space-x-1 text-gray-500 text-xs mt-1">
            <Calendar className="h-3 w-3" />
            <span>{dateFormatted}</span>
          </div>
        </div>

        <span className="px-2.5 py-1 rounded-md text-xs font-bold tracking-wider bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
          {account.currency}
        </span>
      </div>

      {/* Balance Display */}
      <div className="my-5">
        <p className="text-xs text-gray-400 uppercase tracking-wider font-medium">Available Balance</p>
        <p className="text-2xl sm:text-3xl font-black text-white mt-1 tracking-tight">
          <span className="text-gray-400 font-normal mr-1">{symbol}</span>
          {account.balance.toLocaleString()}
        </p>
      </div>

      {/* Quick Action Buttons */}
      <div className="grid grid-cols-2 gap-2 pt-3 border-t border-gray-800/80">
        <button
          onClick={() => onTransferFrom(account)}
          className="flex items-center justify-center space-x-1.5 py-2 px-3 rounded-lg bg-gray-800/80 hover:bg-gray-700 text-xs font-semibold text-gray-200 transition border border-gray-700/50 hover:border-gray-600"
        >
          <ArrowUpRight className="h-3.5 w-3.5 text-indigo-400" />
          <span>Send From</span>
        </button>

        <button
          onClick={() => onTransferTo(account)}
          className="flex items-center justify-center space-x-1.5 py-2 px-3 rounded-lg bg-gray-800/80 hover:bg-gray-700 text-xs font-semibold text-gray-200 transition border border-gray-700/50 hover:border-gray-600"
        >
          <ArrowDownLeft className="h-3.5 w-3.5 text-emerald-400" />
          <span>Receive To</span>
        </button>
      </div>
    </div>
  )
}
