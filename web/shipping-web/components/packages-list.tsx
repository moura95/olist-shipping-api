"use client"

import { CardContent } from "@/components/ui/card"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { useToast } from "@/hooks/use-toast"
import type { Package, Carrier, State } from "@/types"
import { Truck, PackageIcon, Calendar, DollarSign, Plus } from "lucide-react"
import { CreatePackageForm } from "@/components/create-package-form"

const API_BASE_URL = "http://18.231.246.36:8080"

interface PackagesListProps {
  packages: Package[]
  carriers: Carrier[]
  states: State[]
  onPackageUpdated: () => void
  onLoadPackages: () => void
  onPackageCreated: () => void
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

export function PackagesList({
  packages,
  carriers,
  states,
  onPackageUpdated,
  onLoadPackages,
  onPackageCreated,
}: PackagesListProps) {
  const [selectedPackage, setSelectedPackage] = useState<Package | null>(null)
  const [newStatus, setNewStatus] = useState("")
  const [hireData, setHireData] = useState({
    carrier_id: "",
    autoQuote: null as any,
  })
  const [loading, setLoading] = useState(false)
  const [loadingQuote, setLoadingQuote] = useState(false)
  const { toast } = useToast()
  const [detailsPackage, setDetailsPackage] = useState<Package | null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)

  useEffect(() => {
    onLoadPackages()
  }, [])

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
          setHireData({ ...hireData, autoQuote: carrierQuote })
        }
      }
    } catch (error) {
      console.error("Erro ao buscar cotação:", error)
    } finally {
      setLoadingQuote(false)
    }
  }

  const updatePackageStatus = async (packageId: string, status: string) => {
    setLoading(true)
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages/${packageId}/status`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ status }),
      })

      if (response.ok) {
        toast({
          title: "Sucesso!",
          description: "Status atualizado com sucesso.",
        })
        onPackageUpdated()
      } else {
        const error = await response.json()
        toast({
          title: "Erro",
          description: error.message || "Erro ao atualizar status.",
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

  const hireCarrier = async (packageId: string) => {
    if (!hireData.carrier_id || !hireData.autoQuote) {
      toast({
        title: "Erro",
        description: "Por favor, selecione uma transportadora.",
        variant: "destructive",
      })
      return
    }

    setLoading(true)
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages/${packageId}/hire`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          transportadora_id: hireData.carrier_id,
          preco: hireData.autoQuote.preco_estimado.toString(),
          prazo_dias: hireData.autoQuote.prazo_estimado_dias,
        }),
      })

      if (response.ok) {
        toast({
          title: "Sucesso!",
          description: "Transportadora contratada com sucesso.",
        })
        setHireData({ carrier_id: "", autoQuote: null })
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

  const getPackageDetails = async (packageId: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages/${packageId}`)
      if (response.ok) {
        const data = await response.json()
        setDetailsPackage(data.data)
      } else {
        toast({
          title: "Erro",
          description: "Erro ao carregar detalhes do pacote.",
          variant: "destructive",
        })
      }
    } catch (error) {
      toast({
        title: "Erro",
        description: "Erro de conexão com a API.",
        variant: "destructive",
      })
    }
  }

  if (packages.length === 0) {
    return (
      <div className="text-center py-8">
        <PackageIcon className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
        <p className="text-muted-foreground">Nenhum pacote encontrado</p>
        <Button onClick={onLoadPackages} variant="outline" className="mt-4">
          Recarregar
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold">Total: {packages.length} pacotes</h3>
        <div className="flex gap-2">
          <Button onClick={() => setShowCreateModal(true)} size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Criar Pacote
          </Button>
        </div>
      </div>

      <div className="grid gap-4">
        {packages.map((pkg) => (
          <Card
            key={pkg.id}
            className="p-4 cursor-pointer hover:shadow-md transition-shadow"
            onClick={() => getPackageDetails(pkg.id!)}
          >
            <div className="flex justify-between items-start">
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <h4 className="font-semibold">{pkg.produto}</h4>
                  <Badge
                    className={statusColors[pkg.status as keyof typeof statusColors] || "bg-gray-100 text-gray-800"}
                  >
                    {statusLabels[pkg.status as keyof typeof statusLabels] || pkg.status}
                  </Badge>
                </div>

                <div className="text-sm text-muted-foreground space-y-1">
                  <p>
                    <strong>Código:</strong> {pkg.codigo_rastreio}
                  </p>
                  <p>
                    <strong>Peso:</strong> {pkg.peso_kg}kg
                  </p>
                  <p>
                    <strong>Destino:</strong> {pkg.estado_destino}
                  </p>
                  {pkg.transportadora_id && (
                    <div className="flex items-center gap-4 text-green-600">
                      <span className="flex items-center gap-1">
                        <Truck className="h-4 w-4" />
                        Transportadora contratada
                      </span>
                      {pkg.preco_contratado && (
                        <span className="flex items-center gap-1">
                          <DollarSign className="h-4 w-4" />
                          R$ {pkg.preco_contratado}
                        </span>
                      )}
                      {pkg.prazo_contratado_dias && (
                        <span className="flex items-center gap-1">
                          <Calendar className="h-4 w-4" />
                          {pkg.prazo_contratado_dias} dias
                        </span>
                      )}
                    </div>
                  )}
                </div>
              </div>

              <div className="flex gap-2" onClick={(e) => e.stopPropagation()}>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant="outline" size="sm">
                      Atualizar Status
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Atualizar Status do Pacote</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                      <Select value={newStatus} onValueChange={setNewStatus}>
                        <SelectTrigger>
                          <SelectValue placeholder="Selecione o novo status" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="criado">Criado</SelectItem>
                          <SelectItem value="esperando_coleta">Esperando Coleta</SelectItem>
                          <SelectItem value="coletado">Coletado</SelectItem>
                          <SelectItem value="enviado">Enviado</SelectItem>
                          <SelectItem value="entregue">Entregue</SelectItem>
                          <SelectItem value="extraviado">Extraviado</SelectItem>
                        </SelectContent>
                      </Select>
                      <Button
                        onClick={() => {
                          updatePackageStatus(pkg.id!, newStatus)
                          setNewStatus("")
                        }}
                        disabled={!newStatus || loading}
                        className="w-full"
                      >
                        {loading ? "Atualizando..." : "Atualizar Status"}
                      </Button>
                    </div>
                  </DialogContent>
                </Dialog>

                {!pkg.transportadora_id && (
                  <Dialog>

                    <DialogContent>
                      <div className="space-y-4">
                        <div className="space-y-2">
                          <Label>Transportadora</Label>
                          <Select
                            value={hireData.carrier_id}
                            onValueChange={(value) => {
                              setHireData({ ...hireData, carrier_id: value, autoQuote: null })
                              if (value && pkg) {
                                getAutoQuote(pkg, value)
                              }
                            }}
                          >
                            <SelectTrigger>
                              <SelectValue placeholder="Selecione a transportadora" />
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

                        {loadingQuote && (
                          <div className="text-center py-4">
                            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mx-auto mb-2"></div>
                            <p className="text-sm text-muted-foreground">Buscando cotação...</p>
                          </div>
                        )}

                        {hireData.autoQuote && !loadingQuote && (
                          <Card className="border-green-200 bg-green-50">
                            <CardContent className="p-4">
                              <div className="grid grid-cols-2 gap-4">
                                <div>
                                  <Label className="text-sm font-medium text-muted-foreground">Preço</Label>
                                  <p className="text-sm font-semibold text-green-600">
                                    R$ {hireData.autoQuote.preco_estimado?.toFixed(2)}
                                  </p>
                                </div>
                                <div>
                                  <Label className="text-sm font-medium text-muted-foreground">Prazo</Label>
                                  <p className="text-sm flex items-center gap-1">
                                    <Calendar className="h-3 w-3" />
                                    {hireData.autoQuote.prazo_estimado_dias} dias
                                  </p>
                                </div>
                              </div>
                            </CardContent>
                          </Card>
                        )}

                        <Button
                          onClick={() => hireCarrier(pkg.id!)}
                          disabled={!hireData.carrier_id || !hireData.autoQuote || loading || loadingQuote}
                          className="w-full"
                        >
                          {loading ? "Contratando..." : "Contratar"}
                        </Button>
                      </div>
                    </DialogContent>
                  </Dialog>
                )}
              </div>
            </div>
          </Card>
        ))}
      </div>

      <Dialog open={showCreateModal} onOpenChange={setShowCreateModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Criar Novo Pacote</DialogTitle>
          </DialogHeader>
          <CreatePackageForm
            states={states}
            onPackageCreated={() => {
              onPackageCreated()
              setShowCreateModal(false)
            }}
          />
        </DialogContent>
      </Dialog>

      <Dialog open={!!detailsPackage} onOpenChange={() => setDetailsPackage(null)}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Detalhes do Pacote</DialogTitle>
          </DialogHeader>
          {detailsPackage && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-3">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Produto</Label>
                    <p className="text-sm">{detailsPackage.produto}</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Código de Rastreamento</Label>
                    <p className="text-sm font-mono">{detailsPackage.codigo_rastreio}</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Peso</Label>
                    <p className="text-sm">{detailsPackage.peso_kg}kg</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Estado de Destino</Label>
                    <p className="text-sm">{detailsPackage.estado_destino}</p>
                  </div>
                </div>

                <div className="space-y-3">
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Status Atual</Label>
                    <div className="mt-1">
                      <Badge
                        className={
                          statusColors[detailsPackage.status as keyof typeof statusColors] ||
                          "bg-gray-100 text-gray-800"
                        }
                      >
                        {statusLabels[detailsPackage.status as keyof typeof statusLabels] || detailsPackage.status}
                      </Badge>
                    </div>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Criado em</Label>
                    <p className="text-sm">
                      {detailsPackage.criado_em ? new Date(detailsPackage.criado_em).toLocaleString("pt-BR") : "N/A"}
                    </p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Última Atualização</Label>
                    <p className="text-sm">
                      {detailsPackage.atualizado_em
                        ? new Date(detailsPackage.atualizado_em).toLocaleString("pt-BR")
                        : "N/A"}
                    </p>
                  </div>
                </div>
              </div>

              {detailsPackage.transportadora_id && (
                <div className="border-t pt-4">
                  <h4 className="font-semibold mb-3 flex items-center gap-2">
                    <Truck className="h-4 w-4" />
                    Transportadora Contratada
                  </h4>
                  <div className="grid grid-cols-3 gap-4">
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">ID da Transportadora</Label>
                      <p className="text-sm font-mono">{detailsPackage.transportadora_id}</p>
                    </div>
                    {detailsPackage.preco_contratado && (
                      <div>
                        <Label className="text-sm font-medium text-muted-foreground">Preço Contratado</Label>
                        <p className="text-sm font-semibold text-green-600">R$ {detailsPackage.preco_contratado}</p>
                      </div>
                    )}
                    {detailsPackage.prazo_contratado_dias && (
                      <div>
                        <Label className="text-sm font-medium text-muted-foreground">Prazo de Entrega</Label>
                        <p className="text-sm">{detailsPackage.prazo_contratado_dias} dias</p>
                      </div>
                    )}
                  </div>
                </div>
              )}

              <div className="border-t pt-4">
                <h4 className="font-semibold mb-3">Ações Rápidas</h4>
                <div className="flex gap-2">
                  <Select value={newStatus} onValueChange={setNewStatus}>
                    <SelectTrigger className="w-48">
                      <SelectValue placeholder="Alterar status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="criado">Criado</SelectItem>
                      <SelectItem value="esperando_coleta">Esperando Coleta</SelectItem>
                      <SelectItem value="coletado">Coletado</SelectItem>
                      <SelectItem value="enviado">Enviado</SelectItem>
                      <SelectItem value="entregue">Entregue</SelectItem>
                      <SelectItem value="extraviado">Extraviado</SelectItem>
                    </SelectContent>
                  </Select>
                  <Button
                    onClick={() => {
                      updatePackageStatus(detailsPackage.id!, newStatus)
                      setDetailsPackage(null)
                      setNewStatus("")
                    }}
                    disabled={!newStatus || loading}
                    size="sm"
                  >
                    {loading ? "Atualizando..." : "Atualizar"}
                  </Button>
                </div>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
