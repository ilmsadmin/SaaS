'use client'

import { createContext, useContext, ReactNode } from 'react'
import { useToast as useToastHook, Toast } from '@/hooks/use-toast'
import { ToastContainer } from './toast-container'

interface ToastContextType {
  toast: (message: string, type?: Toast['type'], options?: Partial<Toast>) => string
  success: (message: string, options?: Partial<Toast>) => string
  error: (message: string, options?: Partial<Toast>) => string
  warning: (message: string, options?: Partial<Toast>) => string
  info: (message: string, options?: Partial<Toast>) => string
  removeToast: (id: string) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

export function ToastProvider({ children }: { children: ReactNode }) {
  const toastHook = useToastHook()

  return (
    <ToastContext.Provider value={toastHook}>
      {children}
      <ToastContainer toasts={toastHook.toasts} onRemove={toastHook.removeToast} />
    </ToastContext.Provider>
  )
}

export function useToast() {
  const context = useContext(ToastContext)
  if (context === undefined) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}
