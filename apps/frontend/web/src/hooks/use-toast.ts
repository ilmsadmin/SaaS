import { useState, useCallback } from 'react'

export interface Toast {
  id: string
  title?: string
  message: string
  type: 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

let toastCounter = 0

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>([])

  const addToast = useCallback((toast: Omit<Toast, 'id'>) => {
    const id = `toast-${++toastCounter}`
    const newToast: Toast = {
      id,
      duration: 5000,
      ...toast,
    }

    setToasts(prev => [...prev, newToast])

    // Auto remove toast after duration
    if (newToast.duration && newToast.duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, newToast.duration)
    }

    return id
  }, [])

  const removeToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id))
  }, [])

  const toast = useCallback((message: string, type: Toast['type'] = 'info', options?: Partial<Toast>) => {
    return addToast({
      message,
      type,
      ...options,
    })
  }, [addToast])

  const success = useCallback((message: string, options?: Partial<Toast>) => {
    return toast(message, 'success', options)
  }, [toast])

  const error = useCallback((message: string, options?: Partial<Toast>) => {
    return toast(message, 'error', options)
  }, [toast])

  const warning = useCallback((message: string, options?: Partial<Toast>) => {
    return toast(message, 'warning', options)
  }, [toast])

  const info = useCallback((message: string, options?: Partial<Toast>) => {
    return toast(message, 'info', options)
  }, [toast])

  return {
    toasts,
    toast,
    success,
    error,
    warning,
    info,
    removeToast,
  }
}
