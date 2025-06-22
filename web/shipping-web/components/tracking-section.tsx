"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useToast } from "@/hooks/use-toast"
import type { Package, Carrier } from "@/types"
import { Search, PackageIcon, Truck, Calendar, DollarSign, MapPin } from "lucide-react"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080"

interface TrackingSectionProps {
  carriers: Carrier[]
}

const statusColors = {
  criado: "bg-blue-100 text-blue-800",
  esperando_coleta: "bg-yellow-100 text-yellow-800",
  coletado: "bg-orange-100 text-orange-800",
  enviado: "bg-purple-100 text-purple-800",
  entregue: "bg-green-100 text-green-800",
  extraviado: "bg-red-100 text-red-800",
}

const statusLabels = {
  criado: "Criado",
  esperando_coleta: "Esperando Coleta",
  coletado: "Coletado",
  enviado: "Enviado",
  entregue: "Entregue",
  extraviado: "Extraviado",
}

export function TrackingSection({ carriers }: TrackingSectionProps) {
  const [trackingCode, setTrackingCode] = useState("")
  const [packageData, setPackageData] = useState<Package | null>(null)
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)
  const { toast } = useToast()

  const searchPackage = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!trackingCode.trim()) {
      toast({
        title: "Erro",
        description: "Por favor, informe o código de rastreamento.",
        variant: "destructive",
      })
      return
    }

    setLoading(true)
    setSearched(true)

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages/tracking/${trackingCode}`)

      if (response.ok) {
        const data = await response.json()
        setPackageData(data.data)
      } else if (response.status === 404) {
        setPackageData(null)
        toast({
          title: "Pacote não encontrado",
          description: "Nenhum pacote foi encontrado com este código de rastreamento.",
          variant: "destructive",
        })
      } else {
        const error = await response.json()
        setPackageData(null)
        toast({
          title: "Erro",
          description: error.message || "Erro ao buscar pacote.",
          variant: "destructive",
        })
      }
    } catch (error) {
      setPackageData(null)
      toast({
        title: "Erro",
        description: "Erro de conexão com a API.",
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  const getCarrierName = (carrierId: string) => {
    const carrier = carriers.find((c) => c.id === carrierId)
    return carrier?.nome || "Transportadora não identificada"
  }

  return (
    <div className="space-y-6">
      <form onSubmit={searchPackage} className="space-y-4 max-w-md">
        <div className="space-y-2">
          <Label htmlFor="tracking-code">Código de Rastreamento</Label>
          <div className="flex gap-2">
            <Input
              id="tracking-code"
              value={trackingCode}
              onChange={(e) => setTrackingCode(e.target.value)}
              placeholder="Ex: PKG123456789"
              className="flex-1"
            />
            <Button type="submit" disabled={loading}>
              {loading ? (
                "Buscando..."
              ) : (
                <>
                  <Search className="h-4 w-4 mr-2" />
                  Buscar
                </>
              )}
            </Button>
          </div>
        </div>
      </form>

      {searched && !packageData && !loading && (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-8">
            <PackageIcon className="h-12 w-12 text-muted-foreground mb-4" />
            <p className="text-muted-foreground text-center">Nenhum pacote encontrado com o código "{trackingCode}"</p>
          </CardContent>
        </Card>
      )}

      {packageData && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <PackageIcon className="h-5 w-5" />
              Informações do Pacote
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Produto</Label>
                  <p className="text-lg font-semibold">{packageData.produto}</p>
                </div>

                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Código de Rastreamento</Label>
                  <p className="text-sm font-mono bg-muted px-2 py-1 rounded inline-block">
                    {packageData.codigo_rastreio}
                  </p>
                </div>

                <div className="flex gap-4">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Peso</Label>
                    <p className="text-sm">{packageData.peso_kg}kg</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Destino</Label>
                    <p className="text-sm flex items-center gap-1">
                      <MapPin className="h-3 w-3" />
                      {packageData.estado_destino}
                    </p>
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Status Atual</Label>
                  <div className="mt-1">
                    <Badge
                      className={
                        statusColors[packageData.status as keyof typeof statusColors] || "bg-gray-100 text-gray-800"
                      }
                    >
                      {statusLabels[packageData.status as keyof typeof statusLabels] || packageData.status}
                    </Badge>
                  </div>
                </div>

                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Criado em</Label>
                  <p className="text-sm">
                    {packageData.criado_em ? new Date(packageData.criado_em).toLocaleString("pt-BR") : "N/A"}
                  </p>
                </div>

                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Última Atualização</Label>
                  <p className="text-sm">
                    {packageData.atualizado_em ? new Date(packageData.atualizado_em).toLocaleString("pt-BR") : "N/A"}
                  </p>
                </div>
              </div>
            </div>

            {packageData.transportadora_id && (
              <div className="border-t pt-4">
                <h4 className="font-semibold mb-3 flex items-center gap-2 text-green-700">
                  <Truck className="h-4 w-4" />
                  Transportadora Contratada
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 bg-green-50 p-4 rounded-lg">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Transportadora</Label>
                    <p className="text-sm font-semibold">{getCarrierName(packageData.transportadora_id)}</p>
                  </div>
                  {packageData.preco_contratado && (
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Preço Contratado</Label>
                      <p className="text-sm font-semibold text-green-600 flex items-center gap-1">
                        <DollarSign className="h-3 w-3" />
                        R$ {packageData.preco_contratado}
                      </p>
                    </div>
                  )}
                  {packageData.prazo_contratado_dias && (
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Prazo de Entrega</Label>
                      <p className="text-sm flex items-center gap-1">
                        <Calendar className="h-3 w-3" />
                        {packageData.prazo_contratado_dias} dias
                      </p>
                    </div>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
