import React from 'react'

export const QuickControls: React.FC = () => {
  const handleAction = (action: string) => {
    console.log(`Action: ${action}`)
    // Implementar ações aqui
  }

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800 flex items-center gap-2">
        <i className="ri-tools-line text-orange-600"></i>
        CONTROLES RÁPIDOS
      </h2>
      <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <button
            onClick={() => handleAction('restart')}
            className="p-4 rounded-xl border transition-all duration-200 bg-blue-50 hover:bg-blue-100 border-blue-200 hover:shadow-md"
          >
            <div className="flex flex-col items-center text-center space-y-3">
              <div className="w-12 h-12 rounded-xl flex items-center justify-center border bg-blue-50 hover:bg-blue-100 border-blue-200">
                <i className="ri-restart-line text-blue-600 text-2xl"></i>
              </div>
              <div>
                <div className="font-semibold text-sm text-blue-600">Reiniciar Componente</div>
                <div className="text-xs text-gray-600 mt-1">Reinicia componentes específicos do sistema</div>
              </div>
            </div>
          </button>
          <button
            onClick={() => handleAction('export-logs')}
            className="p-4 rounded-xl border transition-all duration-200 bg-green-50 hover:bg-green-100 border-green-200 hover:shadow-md"
          >
            <div className="flex flex-col items-center text-center space-y-3">
              <div className="w-12 h-12 rounded-xl flex items-center justify-center border bg-green-50 hover:bg-green-100 border-green-200">
                <i className="ri-download-line text-green-600 text-2xl"></i>
              </div>
              <div>
                <div className="font-semibold text-sm text-green-600">Exportar Logs</div>
                <div className="text-xs text-gray-600 mt-1">Baixa logs do sistema em formato ZIP</div>
              </div>
            </div>
          </button>
          <button
            onClick={() => handleAction('settings')}
            className="p-4 rounded-xl border transition-all duration-200 bg-purple-50 hover:bg-purple-100 border-purple-200 hover:shadow-md"
          >
            <div className="flex flex-col items-center text-center space-y-3">
              <div className="w-12 h-12 rounded-xl flex items-center justify-center border bg-purple-50 hover:bg-purple-100 border-purple-200">
                <i className="ri-settings-3-line text-purple-600 text-2xl"></i>
              </div>
              <div>
                <div className="font-semibold text-sm text-purple-600">Configurações</div>
                <div className="text-xs text-gray-600 mt-1">Acessa painel de configurações avançadas</div>
              </div>
            </div>
          </button>
          <button
            onClick={() => handleAction('report')}
            className="p-4 rounded-xl border transition-all duration-200 bg-orange-50 hover:bg-orange-100 border-orange-200 hover:shadow-md"
          >
            <div className="flex flex-col items-center text-center space-y-3">
              <div className="w-12 h-12 rounded-xl flex items-center justify-center border bg-orange-50 hover:bg-orange-100 border-orange-200">
                <i className="ri-file-chart-line text-orange-600 text-2xl"></i>
              </div>
              <div>
                <div className="font-semibold text-sm text-orange-600">Relatório</div>
                <div className="text-xs text-gray-600 mt-1">Gera relatório de conformidade e performance</div>
              </div>
            </div>
          </button>
        </div>
      </div>
      <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
        <h3 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
          <i className="ri-flashlight-line text-yellow-600"></i>
          AÇÕES RÁPIDAS
        </h3>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="p-4 bg-gray-50 rounded-xl border border-gray-200">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium text-gray-700">Cache Flush</span>
              <button
                onClick={() => handleAction('cache-flush')}
                className="text-blue-600 hover:text-blue-700"
              >
                <i className="ri-refresh-line"></i>
              </button>
            </div>
            <div className="text-xs text-gray-600">Limpa todos os níveis de cache</div>
          </div>
          <div className="p-4 bg-gray-50 rounded-xl border border-gray-200">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium text-gray-700">Health Check</span>
              <button
                onClick={() => handleAction('health-check')}
                className="text-green-600 hover:text-green-700"
              >
                <i className="ri-heart-pulse-line"></i>
              </button>
            </div>
            <div className="text-xs text-gray-600">Executa verificação completa</div>
          </div>
          <div className="p-4 bg-gray-50 rounded-xl border border-gray-200">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium text-gray-700">Backup</span>
              <button
                onClick={() => handleAction('backup')}
                className="text-purple-600 hover:text-purple-700"
              >
                <i className="ri-save-line"></i>
              </button>
            </div>
            <div className="text-xs text-gray-600">Cria backup do estado atual</div>
          </div>
        </div>
      </div>
    </div>
  )
}

