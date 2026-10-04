import React from 'react'
import type { Account, Transfer } from '../types'
import { Users, Coins, ArrowUpRight, ShieldCheck } from 'lucide-react'

interface StatsOverviewProps {
  accounts: Account[]
  transfers: Transfer[]
}

export const StatsOverview: React.FC<StatsOverviewProps> = ({ accounts, transfers }) => {
  // Aggregate balance by currency
  const currencyTotals = accounts.reduce((acc, account) => {
    acc[account.currency] = (acc[account.currency] || 0) + account.balance
    return acc
  }, {} as Record<string, number>)

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      {/* Total Accounts */}
      <div className="bg-gray-900/60 border border-gray-800 rounded-xl p-5 backdrop-blur">
        <div className="flex items-center justify-between">
          <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Active Accounts</span>
          <div className="p-2 bg-indigo-500/10 rounded-lg text-indigo-400">
            <Users className="h-5 w-5" />
          </div>
        </div>
        <p className="text-3xl font-extrabold text-white mt-2">{accounts.length}</p>
        <p className="text-xs text-gray-500 mt-1">Unique ledger accounts</p>
      </div>

      {/* Primary Circulation (USD) */}
      <div className="bg-gray-900/60 border border-gray-800 rounded-xl p-5 backdrop-blur">
        <div className="flex items-center justify-between">
          <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">USD Total Reserves</span>
          <div className="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
            <Coins className="h-5 w-5" />
          </div>
        </div>
        <p className="text-3xl font-extrabold text-white mt-2">
          ${(currencyTotals['USD'] || 0).toLocaleString()}
        </p>
        <p className="text-xs text-gray-500 mt-1">Across all USD accounts</p>
      </div>

      {/* Transfers processed */}
      <div className="bg-gray-900/60 border border-gray-800 rounded-xl p-5 backdrop-blur">
        <div className="flex items-center justify-between">
          <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Transfers Executed</span>
          <div className="p-2 bg-violet-500/10 rounded-lg text-violet-400">
            <ArrowUpRight className="h-5 w-5" />
          </div>
        </div>
        <p className="text-3xl font-extrabold text-white mt-2">{transfers.length}</p>
        <p className="text-xs text-gray-500 mt-1">Atomic ledger transactions</p>
      </div>

      {/* Ledger Safety Status */}
      <div className="bg-gray-900/60 border border-gray-800 rounded-xl p-5 backdrop-blur">
        <div className="flex items-center justify-between">
          <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Integrity Engine</span>
          <div className="p-2 bg-amber-500/10 rounded-lg text-amber-400">
            <ShieldCheck className="h-5 w-5" />
          </div>
        </div>
        <div className="flex items-center space-x-1.5 mt-2">
          <span className="inline-block w-2 h-2 rounded-full bg-emerald-400"></span>
          <span className="text-lg font-bold text-emerald-400">ACID Enforced</span>
        </div>
        <p className="text-xs text-gray-500 mt-1">Non-negative checks & deadlocks safe</p>
      </div>
    </div>
  )
}
