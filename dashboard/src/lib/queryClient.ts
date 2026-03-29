import { QueryClient } from "@tanstack/react-query"

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5,      // 5 minutes — avoid redundant API calls
      gcTime: 1000 * 60 * 10,         // 10 minutes — keep unused cache a bit longer
      retry: 2,
      refetchOnWindowFocus: false,     // ISP dashboard; no need to refetch on tab focus
    },
    mutations: {
      retry: 0,
    },
  },
})
