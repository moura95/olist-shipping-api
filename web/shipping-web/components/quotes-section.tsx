"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useToast } from "@/hooks/use-toast"
import type { State, Quote } from "@/types"
import { Truck, Clock, DollarSign } from "lucide-react"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080"

interface QuotesSectionProps {
  states: State[]
}

export function QuotesSection({ states }: QuotesSectionProps) {
  const [formData, setFormData] = useState({
    estado_destino: "",
    peso_kg: "",
  })
  const [quotes, setQuotes] = useState<Quote[]>([])
  const [loading, setLoading] = useState(false)
  const { toast } = useToast()

  const getQuotes = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const params = new URLSearchParams({
        estado_destino: formData.estado_destino,
        peso_kg: formData.peso_kg,
      })

      const response = await fetch(`${API_BASE_URL}/api/v1/quotes?${params}`)

      if (response.ok) {
        const data = await response.json()
        setQuotes(data.data || [])

        if (!data.data || data.data.length === 0) {
          toast({
            title: "Aviso",
            description: "Nenhuma cotação encontrada para os parâmetros informados.",
          })
        }
      } else {
        const error = await response.json()
        toast({
          title: "Erro",
          description: error.message || "Erro ao consultar cotações.",
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

  return (
    <div className="space-y-6">
      <form onSubmit={getQuotes} className="space-y-4 max-w-md">
        <div className="space-y-2">
          <Label htmlFor="quote-state">Estado de Destino</Label>
          <Select
            value={formData.estado_destino}
            onValueChange={(value) => setFormData({ ...formData, estado_destino: value })}
            required
          >
            <SelectTrigger>
              <SelectValue placeholder="Selecione o estado" />
            </SelectTrigger>
            <SelectContent>
              {states.map((state) => (
                <SelectItem key={state.codigo} value={state.codigo || ""}>
                  {state.codigo} - {state.nome}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="quote-weight">Peso (kg)</Label>
          <Input
            id="quote-weight"
            type="number"
            step="0.1"
            min="0.1"
            value={formData.peso_kg}
            onChange={(e) => setFormData({ ...formData, peso_kg: e.target.value })}
            placeholder="Ex: 0.6"
            required
          />
        </div>

        <Button type="submit" disabled={loading} className="w-full">
          {loading ? "Consultando..." : "Consultar Cotações"}
        </Button>
      </form>

      {quotes.length > 0 && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold">Cotações Disponíveis</h3>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {quotes.map((quote, index) => (
              <Card key={index}>
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center gap-2 text-base">
                    <Truck className="h-5 w-5" />
                    {quote.transportadora}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex items-center gap-2 text-green-600">
                    <DollarSign className="h-4 w-4" />
                    <span className="font-semibold">R$ {quote.preco_estimado?.toFixed(2)}</span>
                  </div>
                  <div className="flex items-center gap-2 text-blue-600">
                    <Clock className="h-4 w-4" />
                    <span>{quote.prazo_estimado_dias} dias</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
