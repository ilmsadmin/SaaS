'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { verifyEmail, resendEmailVerification } from '@/lib/api/user-management'
import { useToast } from '@/hooks/use-toast'
import Link from 'next/link'

export default function VerifyEmailPage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const { toast } = useToast()
  
  const [isLoading, setIsLoading] = useState(false)
  const [isResending, setIsResending] = useState(false)
  const [token, setToken] = useState<string | null>(null)
  const [email, setEmail] = useState<string | null>(null)
  const [verificationStatus, setVerificationStatus] = useState<'pending' | 'success' | 'error' | 'expired'>('pending')
  const [errorMessage, setErrorMessage] = useState('')

  useEffect(() => {
    const tokenParam = searchParams.get('token')
    const emailParam = searchParams.get('email')
    
    setToken(tokenParam)
    setEmail(emailParam)

    if (tokenParam) {
      handleVerification(tokenParam)
    }
  }, [searchParams])

  const handleVerification = async (verificationToken: string) => {
    setIsLoading(true)
    setVerificationStatus('pending')

    try {
      await verifyEmail(verificationToken)
      setVerificationStatus('success')
      toast('Email verified successfully! You can now access your account.', 'success')
      
      // Redirect to login after a short delay
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    } catch (err: any) {
      setVerificationStatus('error')
      const message = err.response?.data?.error || 'Failed to verify email'
      setErrorMessage(message)
      
      if (message.includes('expired') || message.includes('invalid')) {
        setVerificationStatus('expired')
      }
    } finally {
      setIsLoading(false)
    }
  }

  const handleResendVerification = async () => {
    if (!email) {
      toast('Email address not found. Please try registering again.', 'error')
      return
    }

    setIsResending(true)

    try {
      await resendEmailVerification({ email })
      toast('Verification email sent! Please check your inbox.', 'success')
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to resend verification email'
      toast(errorMessage, 'error')
    } finally {
      setIsResending(false)
    }
  }

  const renderContent = () => {
    switch (verificationStatus) {
      case 'pending':
        return (
          <div className="text-center">
            <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-blue-100">
              {isLoading ? (
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
              ) : (
                <svg
                  className="h-6 w-6 text-blue-600"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                  />
                </svg>
              )}
            </div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              {isLoading ? 'Verifying your email...' : 'Email Verification'}
            </h2>
            <p className="mt-2 text-center text-sm text-gray-600">
              {isLoading 
                ? 'Please wait while we verify your email address.'
                : 'Click the verification link in your email to verify your account.'
              }
            </p>
          </div>
        )

      case 'success':
        return (
          <div className="text-center">
            <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100">
              <svg
                className="h-6 w-6 text-green-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Email Verified!
            </h2>
            <p className="mt-2 text-center text-sm text-gray-600">
              Your email has been successfully verified. You will be redirected to the login page shortly.
            </p>
            <div className="mt-6">
              <Link
                href="/login"
                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Go to Login
              </Link>
            </div>
          </div>
        )

      case 'expired':
        return (
          <div className="text-center">
            <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100">
              <svg
                className="h-6 w-6 text-yellow-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 14.5c-.77.833.192 2.5 1.732 2.5z"
                />
              </svg>
            </div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Link Expired
            </h2>
            <p className="mt-2 text-center text-sm text-gray-600">
              This verification link has expired or is invalid.
            </p>
            {email && (
              <div className="mt-6">
                <button
                  onClick={handleResendVerification}
                  disabled={isResending}
                  className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                >
                  {isResending ? 'Sending...' : 'Resend Verification Email'}
                </button>
              </div>
            )}
            <div className="mt-4">
              <Link
                href="/register"
                className="font-medium text-indigo-600 hover:text-indigo-500"
              >
                Back to registration
              </Link>
            </div>
          </div>
        )

      case 'error':
        return (
          <div className="text-center">
            <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
              <svg
                className="h-6 w-6 text-red-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Verification Failed
            </h2>
            <p className="mt-2 text-center text-sm text-gray-600">
              {errorMessage || 'There was an error verifying your email address.'}
            </p>
            {email && (
              <div className="mt-6">
                <button
                  onClick={handleResendVerification}
                  disabled={isResending}
                  className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                >
                  {isResending ? 'Sending...' : 'Resend Verification Email'}
                </button>
              </div>
            )}
            <div className="mt-4">
              <Link
                href="/register"
                className="font-medium text-indigo-600 hover:text-indigo-500"
              >
                Try registering again
              </Link>
            </div>
          </div>
        )

      default:
        return null
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {renderContent()}
      </div>
    </div>
  )
}
