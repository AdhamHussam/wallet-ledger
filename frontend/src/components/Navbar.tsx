import React from 'react'
import { Wallet, RefreshCw, PlusCircle, ArrowRightLeft } from 'lucide-react'

interface NavbarProps {
  isHealthy: boolean | null
  onRefresh: () => void
  onOpenCreateAccount: () => void
  onOpenTransfer: () => void
  isRefreshing: boolean
}

export const Navbar: React.FC<NavbarProps> = ({
  isHealthy,
  onRefresh,
  onOpenCreateAccount,
  onOpenTransfer,
  isRefreshing,
}) => {
  return (
    <header className="border-b border-gray-800 bg-gray-950/80 backdrop-blur sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo & Service status */}
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-3">
              <div className="h-10 w-10 rounded-xl bg-gradient-to-tr from-indigo-600 to-violet-500 flex items-center justify-center shadow-lg shadow-indigo-500/20">
                <Wallet className="h-5 w-5 text-white" />
              </div>
              <div>
                <span className="font-bold text-lg text-white tracking-tight flex items-center gap-2">
                  Wallet Ledger
                  <span className="text-xs px-2 py-0.5 rounded-full bg-indigo-500/10 text-indigo-400 font-mono border border-indigo-500/20">
                    ACID Core
                  </span>
                </span>
                <p className="text-xs text-gray-400">High-Reliability Financial Ledger</p>
              </div>
            </div>

            {/* Health pill */}
            <div className="hidden sm:flex items-center space-x-2 pl-4 border-l border-gray-800">
              <div
                className={`h-2.5 w-2.5 rounded-full ${
                  isHealthy === true
                    ? 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)] animate-pulse'
                    : isHealthy === false
                    ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.6)]'
                    : 'bg-amber-500'
                }`}
              />
              <span className="text-xs font-medium text-gray-300">
                {isHealthy === true ? 'Backend Connected' : isHealthy === false ? 'Disconnected' : 'Checking...'}
              </span>
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex items-center space-x-3">
            <button
              onClick={onRefresh}
              disabled={isRefreshing}
              className="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-800 transition"
              title="Refresh Accounts"
            >
              <RefreshCw className={`h-4 w-4 ${isRefreshing ? 'animate-spin text-indigo-400' : ''}`} />
            </button>

            <button
              onClick={onOpenTransfer}
              className="inline-flex items-center space-x-2 px-3.5 py-2 rounded-lg bg-gray-800 hover:bg-gray-700 text-sm font-medium text-gray-100 border border-gray-700 transition shadow-sm hover:border-gray-600"
            >
              <ArrowRightLeft className="h-4 w-4 text-indigo-400" />
              <span>Transfer Money</span>
            </button>

            <button
              onClick={onOpenCreateAccount}
              className="inline-flex items-center space-x-2 px-3.5 py-2 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-sm font-semibold text-white shadow-lg shadow-indigo-600/30 transition hover:shadow-indigo-500/40"
            >
              <PlusCircle className="h-4 w-4" />
              <span>New Account</span>
            </button>
          </div>
        </div>
      </div>
    </header>
  )
}
