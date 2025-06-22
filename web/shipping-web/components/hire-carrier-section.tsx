"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useToast } from "@/hooks/use-toast"
import type { Package, Carrier } from "@/types"
import { Truck, PackageIcon, DollarSign, Calendar, MapPin, AlertCircle } from "lucide-react"

const API_BASE_URL = "http://18.231.106.0:8080"

interface HireCarrierSectionProps {
  packages: Package[]
  carriers: Carrier[]
  onPackageUpdated: () => void
  onLoadPackages: () => void
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

export function HireCarrierSection({ packages, carriers, onPackageUpdated, onLoadPackages }: HireCarrierSectionProps) {
  const [selectedPackageId, setSelectedPackageId] = useState("")
  const [selectedCarrierId, setSelectedCarrierId] = useState("")
  const [loading, setLoading] = useState(false)
  const { toast } = useToast()
  const [autoQuote, setAutoQuote] = useState<any>(null)
  const [loadingQuote, setLoadingQuote] = useState(false)

  useEffect(() => {
    onLoadPackages()
  }, [])

  // Filtrar apenas pacotes que ainda não têm transportadora contratada
  const availablePackages = packages.filter((pkg) => !pkg.transportadora_id)

  const selectedPackage = packages.find((pkg) => pkg.id === selectedPackageId)

  const getAutoQuote = async (packageData: Package, carrierId: string) => {
    if (!packageData.peso_kg || !packageData.estado_destino) return

    setLoadingQuote(true)
    try {
      const params = new URLSearchParams({
        estado_destino: packageData.estado_destino,
        peso_kg: packageData.peso_kg.toString(),
      })

      const response = await fetch(`${API_BASE_URL}/api/v1/quotes?${params}`)
      if (response.ok) {
        const data = await response.json()
        const quotes = data.data || []
        const carrierQuote = quotes.find((q: any) =>
          carriers.find((c) => c.id === carrierId && c.nome === q.transportadora),
        )

        if (carrierQuote) {
          setAutoQuote(carrierQuote)
        }
      }
    } catch (error) {
      console.error("Erro ao buscar cotação:", error)
    } finally {
      setLoadingQuote(false)
    }
  }

  useEffect(() => {
    if (selectedPackage && selectedCarrierId) {
      getAutoQuote(selectedPackage, selectedCarrierId)
    } else {
      setAutoQuote(null)
    }
  }, [selectedPackageId, selectedCarrierId])

  const hireCarrier = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!selectedPackageId || !selectedCarrierId || !autoQuote) {
      toast({
        title: "Erro",
        description: "Por favor, selecione um pacote e uma transportadora.",
        variant: "destructive",
      })
      return
    }

    setLoading(true)

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages/${selectedPackageId}/hire`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          transportadora_id: selectedCarrierId,
          preco: autoQuote.preco_estimado.toString(),
          prazo_dias: autoQuote.prazo_estimado_dias,
        }),
      })

      if (response.ok) {
        toast({
          title: "Sucesso!",
          description: "Transportadora contratada com sucesso.",
        })
        // Limpar formulário
        setSelectedPackageId("")
        setSelectedCarrierId("")
        setAutoQuote(null)
        onPackageUpdated()
      } else {
        const error = await response.json()
        toast({
          title: "Erro",
          description: error.message || "Erro ao contratar transportadora.",
          variant: "destructive",
        })
      }
    } catch (error) {
      toast({
        title: "Erro",
        description: "Erro de conexão com a API.",
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  if (availablePackages.length === 0) {
    return (
      <div className="text-center py-8">
        <AlertCircle className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold mb-2">Nenhum pacote disponível</h3>
        <p className="text-muted-foreground mb-4">
          Todos os pacotes já possuem transportadora contratada ou não há pacotes cadastrados.
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <form onSubmit={hireCarrier} className="space-y-6 max-w-2xl">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="package-select">Selecionar Pacote</Label>
            <Select value={selectedPackageId} onValueChange={setSelectedPackageId} required>
              <SelectTrigger>
                <SelectValue placeholder="Escolha um pacote" />
              </SelectTrigger>
              <SelectContent>
                {availablePackages.map((pkg) => (
                  <SelectItem key={pkg.id} value={pkg.id!}>
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-xs">{pkg.codigo_rastreio}</span>
                      <span>-</span>
                      <span>{pkg.produto}</span>
                      <span className="text-muted-foreground">({pkg.peso_kg}kg)</span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="carrier-select">Selecionar Transportadora</Label>
            <Select value={selectedCarrierId} onValueChange={setSelectedCarrierId} required>
              <SelectTrigger>
                <SelectValue placeholder="Escolha uma transportadora" />
              </SelectTrigger>
              <SelectContent>
                {carriers.map((carrier) => (
                  <SelectItem key={carrier.id} value={carrier.id!}>
                    {carrier.nome}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {autoQuote && (
          <Button type="submit" disabled={loading || loadingQuote} className="w-full">
            {loading ? (
              "Contratando..."
            ) : (
              <>
                <Truck className="h-4 w-4 mr-2" />
                Contratar Transportadora
              </>
            )}
          </Button>
        )}
      </form>

      {selectedPackage && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <PackageIcon className="h-5 w-5" />
              Detalhes do Pacote Selecionado
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-3">
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Produto</Label>
                  <p className="text-sm font-semibold">{selectedPackage.produto}</p>
                </div>
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Código de Rastreamento</Label>
                  <p className="text-sm font-mono">{selectedPackage.codigo_rastreio}</p>
                </div>
              </div>
              <div className="space-y-3">
                <div className="flex gap-4">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Peso</Label>
                    <p className="text-sm">{selectedPackage.peso_kg}kg</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Destino</Label>
                    <p className="text-sm flex items-center gap-1">
                      <MapPin className="h-3 w-3" />
                      {selectedPackage.estado_destino}
                    </p>
                  </div>
                </div>
                <div>
                  <Label className="text-sm font-medium text-muted-foreground">Status</Label>
                  <div className="mt-1">
                    <Badge
                      className={
                        statusColors[selectedPackage.status as keyof typeof statusColors] || "bg-gray-100 text-gray-800"
                      }
                    >
                      {statusLabels[selectedPackage.status as keyof typeof statusLabels] || selectedPackage.status}
                    </Badge>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {loadingQuote && (
        <Card className="border-blue-200 bg-blue-50">
          <CardContent className="flex items-center justify-center py-8">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p className="text-blue-700">Buscando cotação...</p>
            </div>
          </CardContent>
        </Card>
      )}

      {autoQuote && !loadingQuote && (
        <Card className="border-green-200 bg-green-50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-green-700">
              <DollarSign className="h-5 w-5" />
              Resumo da Contratação
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <Label className="text-sm font-medium text-muted-foreground">Transportadora</Label>
                <p className="text-sm font-semibold">
                  {carriers.find((c) => c.id === selectedCarrierId)?.nome || "N/A"}
                </p>
              </div>
              <div>
                <Label className="text-sm font-medium text-muted-foreground">Preço</Label>
                <p className="text-sm font-semibold text-green-600">R$ {autoQuote.preco_estimado?.toFixed(2)}</p>
              </div>
              <div>
                <Label className="text-sm font-medium text-muted-foreground">Prazo</Label>
                <p className="text-sm flex items-center gap-1">
                  <Calendar className="h-3 w-3" />
                  {autoQuote.prazo_estimado_dias} dias
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
