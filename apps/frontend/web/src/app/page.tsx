export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-6">
            Welcome to Zplus SaaS Platform
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Comprehensive multi-tenant SaaS platform with CRM, HRM, POS, LMS and more business modules
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mt-12">
            <ModuleCard
              title="CRM"
              description="Customer Relationship Management"
              status="available"
              icon="👥"
            />
            <ModuleCard
              title="HRM"
              description="Human Resource Management"
              status="available"
              icon="🏢"
            />
            <ModuleCard
              title="POS"
              description="Point of Sale System"
              status="available"
              icon="🛒"
            />
            <ModuleCard
              title="LMS"
              description="Learning Management System"
              status="available"
              icon="📚"
            />
            <ModuleCard
              title="Check-in"
              description="Attendance Tracking"
              status="available"
              icon="⏰"
            />
            <ModuleCard
              title="Payment"
              description="Payment Processing"
              status="available"
              icon="💳"
            />
            <ModuleCard
              title="Accounting"
              description="Financial Management"
              status="development"
              icon="💰"
            />
            <ModuleCard
              title="E-commerce"
              description="Online Store Platform"
              status="planned"
              icon="🛍️"
            />
          </div>

          <div className="mt-16 flex flex-col sm:flex-row gap-4 justify-center">
            <a
              href="/login"
              className="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Get Started
            </a>
            <a
              href="/demo"
              className="inline-flex items-center px-6 py-3 border border-gray-300 text-base font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              View Demo
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

interface ModuleCardProps {
  title: string
  description: string
  status: 'available' | 'development' | 'planned'
  icon: string
}

function ModuleCard({ title, description, status, icon }: ModuleCardProps) {
  const statusColors = {
    available: 'bg-green-100 text-green-800',
    development: 'bg-yellow-100 text-yellow-800',
    planned: 'bg-gray-100 text-gray-800',
  }

  const statusText = {
    available: 'Available',
    development: 'In Development',
    planned: 'Planned',
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
      <div className="text-3xl mb-4">{icon}</div>
      <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
      <p className="text-gray-600 text-sm mb-4">{description}</p>
      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${statusColors[status]}`}>
        {statusText[status]}
      </span>
    </div>
  )
}
