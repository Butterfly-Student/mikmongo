import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import { createCashEntrySchema, type CreateCashEntry } from '@/lib/schemas/billing'
import { useCreateCashEntry } from '@/hooks/use-cash'
import { cashEntrySources } from '../data/schema'

interface CreateCashEntryDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function CreateCashEntryDialog({ open, onOpenChange }: CreateCashEntryDialogProps) {
  const { mutateAsync: createEntry, isPending } = useCreateCashEntry()

  const today = new Date().toISOString().split('T')[0]

  const form = useForm<CreateCashEntry>({
    resolver: zodResolver(createCashEntrySchema) as never,
    defaultValues: {
      type: 'income',
      source: 'other',
      amount: 0,
      description: '',
      payment_method: '',
      bank_name: '',
      account_number: '',
      entry_date: today,
      notes: '',
    },
  })

  const onSubmit = async (data: CreateCashEntry) => {
    try {
      const payload: CreateCashEntry = {
        type: data.type,
        source: data.source,
        amount: Number(data.amount),
        description: data.description,
        payment_method: data.payment_method,
      }
      if (data.bank_name) payload.bank_name = data.bank_name
      if (data.account_number) payload.account_number = data.account_number
      if (data.entry_date) payload.entry_date = data.entry_date
      if (data.notes) payload.notes = data.notes

      await createEntry(payload)
      form.reset()
      onOpenChange(false)
    } catch {
      // Error handled by hook
    }
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) form.reset()
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className='sm:max-w-[540px] max-h-[90vh] overflow-y-auto'>
        <DialogHeader>
          <DialogTitle>Tambah Entri Kas</DialogTitle>
          <DialogDescription>
            Buat entri kas baru. Isi semua kolom yang diperlukan.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-4'>
            <div className='grid grid-cols-2 gap-4'>
              <FormField
                control={form.control}
                name='type'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tipe *</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      disabled={isPending}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Pilih tipe' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='income'>Masuk</SelectItem>
                        <SelectItem value='expense'>Keluar</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='source'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Sumber *</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      disabled={isPending}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder='Pilih sumber' />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {cashEntrySources.map((s) => (
                          <SelectItem key={s.value} value={s.value}>
                            {s.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='amount'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Jumlah (Rp) *</FormLabel>
                    <FormControl>
                      <Input
                        type='number'
                        min={0}
                        placeholder='0'
                        {...field}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='payment_method'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Metode Pembayaran *</FormLabel>
                    <FormControl>
                      <Input
                        placeholder='Tunai, Transfer, dll'
                        {...field}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem className='col-span-2'>
                    <FormLabel>Deskripsi *</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder='Masukkan deskripsi entri kas'
                        rows={2}
                        {...field}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='bank_name'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nama Bank</FormLabel>
                    <FormControl>
                      <Input
                        placeholder='BCA, Mandiri, dll'
                        {...field}
                        value={field.value ?? ''}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='account_number'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>No. Rekening</FormLabel>
                    <FormControl>
                      <Input
                        placeholder='1234567890'
                        {...field}
                        value={field.value ?? ''}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='entry_date'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tanggal</FormLabel>
                    <FormControl>
                      <Input
                        type='date'
                        {...field}
                        value={field.value ?? ''}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='notes'
                render={({ field }) => (
                  <FormItem className='col-span-2'>
                    <FormLabel>Catatan</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder='Catatan tambahan (opsional)'
                        rows={2}
                        {...field}
                        value={field.value ?? ''}
                        disabled={isPending}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter>
              <Button
                type='button'
                variant='ghost'
                onClick={() => handleOpenChange(false)}
                disabled={isPending}
              >
                Batal
              </Button>
              <Button type='submit' disabled={isPending}>
                {isPending && <Loader2 className='mr-2 size-4 animate-spin' />}
                Simpan
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
