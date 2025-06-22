"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PackagesList } from "@/components/packages-list"
import { QuotesSection } from "@/components/quotes-section"
import { TrackingSection } from "@/components/tracking-section"
import { HireCarrierSection } from "@/components/hire-carrier-section"
import type { Package, Carrier, State } from "@/types"

const API_BASE_URL = "http://18.231.246.36:8080"

export default function Home() {
  const [packages, setPackages] = useState<Package[]>([])
  const [carriers, setCarriers] = useState<Carrier[]>([])
  const [states, setStates] = useState<State[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadInitialData()
  }, [])

  const loadInitialData = async () => {
    try {
      const [carriersRes, statesRes] = await Promise.all([
        fetch(`${API_BASE_URL}/api/v1/carriers`),
        fetch(`${API_BASE_URL}/api/v1/states`),
      ])

      if (carriersRes.ok) {
        const carriersData = await carriersRes.json()
        setCarriers(carriersData.data || [])
      }

      if (statesRes.ok) {
        const statesData = await statesRes.json()
        setStates(statesData.data || [])
      }
    } catch (error) {
      console.error("Erro ao carregar dados iniciais:", error)
    } finally {
      setLoading(false)
    }
  }

  const loadPackages = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages`)
      if (response.ok) {
        const data = await response.json()
        setPackages(data.data || [])
      }
    } catch (error) {
      console.error("Erro ao carregar pacotes:", error)
    }
  }

  const handlePackageCreated = () => {
    loadPackages()
  }

  const handlePackageUpdated = () => {
    loadPackages()
  }

  if (loading) {
    return (
      <div className="container mx-auto p-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-lg">Carregando...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Sistema de Gerenciamento de Fretes</h1>
        <p className="text-muted-foreground">Gerencie pacotes, consulte cotações e contrate transportadoras</p>
      </div>

      <Tabs defaultValue="packages" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="packages">Pacotes</TabsTrigger>
          <TabsTrigger value="tracking">Rastrear Pacote</TabsTrigger>
          <TabsTrigger value="quotes">Cotação de Frete</TabsTrigger>
          <TabsTrigger value="hire">Contratar Transportadora</TabsTrigger>
        </TabsList>

        <TabsContent value="packages">
          <Card>
            <CardHeader>
              <CardTitle>Lista de Pacotes</CardTitle>
              <CardDescription>Visualize e gerencie todos os pacotes cadastrados</CardDescription>
            </CardHeader>
            <CardContent>
              <PackagesList
                packages={packages}
                carriers={carriers}
                states={states}
                onPackageUpdated={handlePackageUpdated}
                onLoadPackages={loadPackages}
                onPackageCreated={handlePackageCreated}
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="tracking">
          <Card>
            <CardHeader>
              <CardTitle>Rastrear Pacote</CardTitle>
              <CardDescription>
                Consulte informações de um pacote específico pelo código de rastreamento
              </CardDescription>
            </CardHeader>
            <CardContent>
              <TrackingSection carriers={carriers} />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="quotes">
          <Card>
            <CardHeader>
              <CardTitle>Cotação de Frete</CardTitle>
              <CardDescription>Consulte preços e prazos de frete de diferentes transportadoras</CardDescription>
            </CardHeader>
            <CardContent>
              <QuotesSection states={states} />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="hire">
          <Card>
            <CardHeader>
              <CardTitle>Contratar Transportadora</CardTitle>
              <CardDescription>
                Selecione um pacote e contrate uma transportadora para realizar a entrega
              </CardDescription>
            </CardHeader>
            <CardContent>
              <HireCarrierSection
                packages={packages}
                carriers={carriers}
                onPackageUpdated={handlePackageUpdated}
                onLoadPackages={loadPackages}
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
