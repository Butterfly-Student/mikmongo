import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface RouterStoreState {
    selectedRouterId: string | null
    selectedRouterName: string | null
    isHydrated: boolean
    setSelectedRouter: (id: string, name: string) => void
    clearSelectedRouter: () => void
    setHydrated: () => void
}

export const useRouterStore = create<RouterStoreState>()(
    persist(
        (set) => ({
            selectedRouterId: null,
            selectedRouterName: null,
            isHydrated: false,
            setSelectedRouter: (id, name) =>
                set((state) => ({
                    ...state,
                    selectedRouterId: id,
                    selectedRouterName: name,
                })),
            clearSelectedRouter: () =>
                set((state) => ({
                    ...state,
                    selectedRouterId: null,
                    selectedRouterName: null,
                })),
            setHydrated: () => set((state) => ({ ...state, isHydrated: true })),
        }),
        {
            name: 'mikmongo-router',
            partialize: (state) => ({
                selectedRouterId: state.selectedRouterId,
                selectedRouterName: state.selectedRouterName,
            }),
            onRehydrateStorage: () => (state) => {
                state?.setHydrated()
            },
        }
    )
)
