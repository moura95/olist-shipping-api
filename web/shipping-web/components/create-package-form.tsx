"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useToast } from "@/hooks/use-toast"
import type { State } from "@/types"


const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "localhost:8080"

interface CreatePackageFormProps {
  states: State[]
  onPackageCreated: () => void
}

export function CreatePackageForm({ states, onPackageCreated }: CreatePackageFormProps) {
  const [formData, setFormData] = useState({
    produto: "",
    peso_kg: "",
    estado_destino: "",
  })
  const [loading, setLoading] = useState(false)
  const { toast } = useToast()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/packages`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          produto: formData.produto,
          peso_kg: Number.parseFloat(formData.peso_kg),
          estado_destino: formData.estado_destino,
        }),
      })

      if (response.ok) {
        toast({
          title: "Sucesso!",
          description: "Pacote criado com sucesso.",
        })
        setFormData({ produto: "", peso_kg: "", estado_destino: "" })
        onPackageCreated()
      } else {
        const error = await response.json()
        toast({
          title: "Erro",
          description: error.message || "Erro ao criar pacote.",
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
    <form onSubmit={handleSubmit} className="space-y-4 max-w-md">
      <div className="space-y-2">
        <Label htmlFor="product">Produto</Label>
        <Input
          id="product"
          value={formData.produto}
          onChange={(e) => setFormData({ ...formData, produto: e.target.value })}
          placeholder="Ex: Camisa tamanho G"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="weight">Peso (kg)</Label>
        <Input
          id="weight"
          type="number"
          step="0.1"
          min="0.1"
          value={formData.peso_kg}
          onChange={(e) => setFormData({ ...formData, peso_kg: e.target.value })}
          placeholder="Ex: 0.6"
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="state">Estado de Destino</Label>
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

      <Button type="submit" disabled={loading} className="w-full">
        {loading ? "Criando..." : "Criar Pacote"}
      </Button>
    </form>
  )
}
