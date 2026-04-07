import { toast } from 'sonner'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { useDeleteUser } from '@/hooks/use-users'
import type { UserResponse } from '@/api/types'
import { useAuthStore } from '@/stores/auth-store'
import { Button } from '@/components/ui/button'

interface DeleteUserDialogProps {
  user: UserResponse | null
  onClose: () => void
}

export function DeleteUserDialog({ user, onClose }: DeleteUserDialogProps) {
  const { mutateAsync: deleteUser, isPending } = useDeleteUser()
  const { adminUser } = useAuthStore()

  if (!user) return null

  const isSelf = adminUser?.id === user.id

  const handleDelete = async () => {
    if (isSelf) return
    try {
      await deleteUser(user.id)
      toast.success(`${user.full_name} has been deleted`)
      onClose()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to delete user')
    }
  }

  return (
    <AlertDialog open={!!user} onOpenChange={(open) => !open && onClose()}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete User</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to permanently delete <span className='font-medium text-foreground'>{user.full_name}</span>? 
            This action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
          {isSelf ? (
            <Button variant='destructive' disabled>
              You cannot delete your own account
            </Button>
          ) : (
            <Button
              variant='destructive'
              onClick={handleDelete}
              disabled={isPending}
            >
              {isPending ? 'Deleting...' : 'Delete'}
            </Button>
          )}
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
