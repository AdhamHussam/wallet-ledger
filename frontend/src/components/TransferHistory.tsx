import React from 'react'
import type { Transfer, Account } from '../types'
import { ArrowUpRight, Clock } from 'lucide-react'

interface TransferHistoryProps {
  transfers: Transfer[]
  accounts: Account[]
}

export const TransferHistory: React.FC<TransferHistoryProps> = ({ transfers, accounts }) => {
  const accountMap = new Map(accounts.map((a) => [a.id, a]))

  if (transfers.length === 0) {
    return (
      <div className="bg-gray-900/40 border border-gray-800 rounded-2xl p-8 text-center">
        <Clock className="h-8 w-8 text-gray-600 mx-auto mb-2" />
        <p className="text-gray-400 font-medium text-sm">No transfers recorded in this session</p>
        <p className="text-gray-600 text-xs mt-1">Transfers executed will appear here in real-time</p>
      </div>
    )
  }

  return (
    <div className="bg-gray-900/60 border border-gray-800 rounded-2xl overflow-hidden shadow-xl">
      <div className="px-6 py-4 border-b border-gray-800 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Clock className="h-4 w-4 text-indigo-400" />
          <h3 className="text-sm font-bold text-white tracking-tight">Recent Transfers & Audit Trail</h3>
        </div>
        <span className="text-xs text-gray-500 font-mono">{transfers.length} records</span>
      </div>

      <div className="divide-y divide-gray-800/60 max-h-96 overflow-y-auto">
        {transfers.map((tx) => {
          const from = accountMap.get(tx.from_account_id)
          const to = accountMap.get(tx.to_account_id)
          const currency = from?.currency || to?.currency || 'USD'
          const dateStr = new Date(tx.created_at).toLocaleTimeString(undefined, {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
          })

          return (
            <div
              key={tx.id}
              className="p-4 hover:bg-gray-800/30 transition flex items-center justify-between text-sm"
            >
              <div className="flex items-center space-x-3">
                <div className="h-9 w-9 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400">
                  <ArrowUpRight className="h-4 w-4" />
                </div>
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="font-semibold text-white">
                      {from?.owner || `Account #${tx.from_account_id}`}
                    </span>
                    <span className="text-gray-500 text-xs">→</span>
                    <span className="font-semibold text-white">
                      {to?.owner || `Account #${tx.to_account_id}`}
                    </span>
                  </div>
                  <div className="flex items-center space-x-2 text-xs text-gray-500 mt-0.5">
                    <span>Tx #{tx.id}</span>
                    <span>•</span>
                    <span>{dateStr}</span>
                  </div>
                </div>
              </div>

              <div className="text-right">
                <p className="font-mono font-bold text-white text-base">
                  {tx.amount.toLocaleString()} {currency}
                </p>
                <span className="inline-flex items-center text-[10px] uppercase font-semibold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded-full border border-emerald-500/20">
                  Committed
                </span>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
