export interface Package {
  id?: string
  codigo_rastreio?: string
  produto?: string
  peso_kg?: number
  estado_destino?: string
  status?: string
  transportadora_id?: string
  preco_contratado?: string
  prazo_contratado_dias?: number
  criado_em?: string
  atualizado_em?: string
}

export interface Quote {
  transportadora?: string
  preco_estimado?: number
  prazo_estimado_dias?: number
}

export interface Carrier {
  id?: string
  nome?: string
  criado_em?: string
}

export interface State {
  codigo?: string
  nome?: string
  nome_regiao?: string
}
