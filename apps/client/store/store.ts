import { configureStore } from '@reduxjs/toolkit'
import authReducer from './slices/auth.slice'

// Define AppStore type first, as 'store' will use it.
// For now, we can define it based on the return type of configureStore directly,
// or use a placeholder and then refine it.
// Let's define it based on the expected return type of configureStore.
// This avoids a circular dependency with `makeStore` if `AppStore` is used in `store`'s declaration.
type StoreType = ReturnType<typeof configureStore<any, any, any>> // A temporary type for the store instance

let store: StoreType | undefined

export const makeStore = () => {
  if (!store) {
    // Ensure store is only initialized once if makeStore is called multiple times
    store = configureStore({
      reducer: {
        auth: authReducer
      },
      devTools: process.env.NODE_ENV !== 'production'
    })
  }
  return store
}

export const getStore = () => store

export type AppStore = ReturnType<typeof makeStore>
export type RootState = ReturnType<AppStore['getState']>
export type AppDispatch = AppStore['dispatch']
