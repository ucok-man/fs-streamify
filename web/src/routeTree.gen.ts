/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file was automatically generated by TanStack Router.
// You should NOT make any changes in this file as it will be overwritten.
// Additionally, you should also exclude this file from your linter and/or formatter to prevent it from being checked or modified.

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as ProtectedRouteImport } from './routes/_protected/route'
import { Route as ProtectedLayoutRouteImport } from './routes/_protected/_layout/route'
import { Route as AuthSignupRouteImport } from './routes/_auth/signup/route'
import { Route as AuthSigninRouteImport } from './routes/_auth/signin/route'
import { Route as AuthOnboardingRouteImport } from './routes/_auth/onboarding/route'
import { Route as ProtectedLayoutIndexImport } from './routes/_protected/_layout/index'
import { Route as ProtectedCallIdRouteImport } from './routes/_protected/call/$id/route'
import { Route as ProtectedLayoutNotificationRouteImport } from './routes/_protected/_layout/notification/route'
import { Route as ProtectedLayoutFriendRouteImport } from './routes/_protected/_layout/friend/route'
import { Route as ProtectedLayoutChatIdRouteImport } from './routes/_protected/_layout/chat/$id/route'

// Create/Update Routes

const ProtectedRouteRoute = ProtectedRouteImport.update({
  id: '/_protected',
  getParentRoute: () => rootRoute,
} as any)

const ProtectedLayoutRouteRoute = ProtectedLayoutRouteImport.update({
  id: '/_layout',
  getParentRoute: () => ProtectedRouteRoute,
} as any)

const AuthSignupRouteRoute = AuthSignupRouteImport.update({
  id: '/_auth/signup',
  path: '/signup',
  getParentRoute: () => rootRoute,
} as any)

const AuthSigninRouteRoute = AuthSigninRouteImport.update({
  id: '/_auth/signin',
  path: '/signin',
  getParentRoute: () => rootRoute,
} as any)

const AuthOnboardingRouteRoute = AuthOnboardingRouteImport.update({
  id: '/_auth/onboarding',
  path: '/onboarding',
  getParentRoute: () => rootRoute,
} as any)

const ProtectedLayoutIndexRoute = ProtectedLayoutIndexImport.update({
  id: '/',
  path: '/',
  getParentRoute: () => ProtectedLayoutRouteRoute,
} as any)

const ProtectedCallIdRouteRoute = ProtectedCallIdRouteImport.update({
  id: '/call/$id',
  path: '/call/$id',
  getParentRoute: () => ProtectedRouteRoute,
} as any)

const ProtectedLayoutNotificationRouteRoute =
  ProtectedLayoutNotificationRouteImport.update({
    id: '/notification',
    path: '/notification',
    getParentRoute: () => ProtectedLayoutRouteRoute,
  } as any)

const ProtectedLayoutFriendRouteRoute = ProtectedLayoutFriendRouteImport.update(
  {
    id: '/friend',
    path: '/friend',
    getParentRoute: () => ProtectedLayoutRouteRoute,
  } as any,
)

const ProtectedLayoutChatIdRouteRoute = ProtectedLayoutChatIdRouteImport.update(
  {
    id: '/chat/$id',
    path: '/chat/$id',
    getParentRoute: () => ProtectedLayoutRouteRoute,
  } as any,
)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/_protected': {
      id: '/_protected'
      path: ''
      fullPath: ''
      preLoaderRoute: typeof ProtectedRouteImport
      parentRoute: typeof rootRoute
    }
    '/_auth/onboarding': {
      id: '/_auth/onboarding'
      path: '/onboarding'
      fullPath: '/onboarding'
      preLoaderRoute: typeof AuthOnboardingRouteImport
      parentRoute: typeof rootRoute
    }
    '/_auth/signin': {
      id: '/_auth/signin'
      path: '/signin'
      fullPath: '/signin'
      preLoaderRoute: typeof AuthSigninRouteImport
      parentRoute: typeof rootRoute
    }
    '/_auth/signup': {
      id: '/_auth/signup'
      path: '/signup'
      fullPath: '/signup'
      preLoaderRoute: typeof AuthSignupRouteImport
      parentRoute: typeof rootRoute
    }
    '/_protected/_layout': {
      id: '/_protected/_layout'
      path: ''
      fullPath: ''
      preLoaderRoute: typeof ProtectedLayoutRouteImport
      parentRoute: typeof ProtectedRouteImport
    }
    '/_protected/_layout/friend': {
      id: '/_protected/_layout/friend'
      path: '/friend'
      fullPath: '/friend'
      preLoaderRoute: typeof ProtectedLayoutFriendRouteImport
      parentRoute: typeof ProtectedLayoutRouteImport
    }
    '/_protected/_layout/notification': {
      id: '/_protected/_layout/notification'
      path: '/notification'
      fullPath: '/notification'
      preLoaderRoute: typeof ProtectedLayoutNotificationRouteImport
      parentRoute: typeof ProtectedLayoutRouteImport
    }
    '/_protected/call/$id': {
      id: '/_protected/call/$id'
      path: '/call/$id'
      fullPath: '/call/$id'
      preLoaderRoute: typeof ProtectedCallIdRouteImport
      parentRoute: typeof ProtectedRouteImport
    }
    '/_protected/_layout/': {
      id: '/_protected/_layout/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof ProtectedLayoutIndexImport
      parentRoute: typeof ProtectedLayoutRouteImport
    }
    '/_protected/_layout/chat/$id': {
      id: '/_protected/_layout/chat/$id'
      path: '/chat/$id'
      fullPath: '/chat/$id'
      preLoaderRoute: typeof ProtectedLayoutChatIdRouteImport
      parentRoute: typeof ProtectedLayoutRouteImport
    }
  }
}

// Create and export the route tree

interface ProtectedLayoutRouteRouteChildren {
  ProtectedLayoutFriendRouteRoute: typeof ProtectedLayoutFriendRouteRoute
  ProtectedLayoutNotificationRouteRoute: typeof ProtectedLayoutNotificationRouteRoute
  ProtectedLayoutIndexRoute: typeof ProtectedLayoutIndexRoute
  ProtectedLayoutChatIdRouteRoute: typeof ProtectedLayoutChatIdRouteRoute
}

const ProtectedLayoutRouteRouteChildren: ProtectedLayoutRouteRouteChildren = {
  ProtectedLayoutFriendRouteRoute: ProtectedLayoutFriendRouteRoute,
  ProtectedLayoutNotificationRouteRoute: ProtectedLayoutNotificationRouteRoute,
  ProtectedLayoutIndexRoute: ProtectedLayoutIndexRoute,
  ProtectedLayoutChatIdRouteRoute: ProtectedLayoutChatIdRouteRoute,
}

const ProtectedLayoutRouteRouteWithChildren =
  ProtectedLayoutRouteRoute._addFileChildren(ProtectedLayoutRouteRouteChildren)

interface ProtectedRouteRouteChildren {
  ProtectedLayoutRouteRoute: typeof ProtectedLayoutRouteRouteWithChildren
  ProtectedCallIdRouteRoute: typeof ProtectedCallIdRouteRoute
}

const ProtectedRouteRouteChildren: ProtectedRouteRouteChildren = {
  ProtectedLayoutRouteRoute: ProtectedLayoutRouteRouteWithChildren,
  ProtectedCallIdRouteRoute: ProtectedCallIdRouteRoute,
}

const ProtectedRouteRouteWithChildren = ProtectedRouteRoute._addFileChildren(
  ProtectedRouteRouteChildren,
)

export interface FileRoutesByFullPath {
  '': typeof ProtectedLayoutRouteRouteWithChildren
  '/onboarding': typeof AuthOnboardingRouteRoute
  '/signin': typeof AuthSigninRouteRoute
  '/signup': typeof AuthSignupRouteRoute
  '/friend': typeof ProtectedLayoutFriendRouteRoute
  '/notification': typeof ProtectedLayoutNotificationRouteRoute
  '/call/$id': typeof ProtectedCallIdRouteRoute
  '/': typeof ProtectedLayoutIndexRoute
  '/chat/$id': typeof ProtectedLayoutChatIdRouteRoute
}

export interface FileRoutesByTo {
  '': typeof ProtectedRouteRouteWithChildren
  '/onboarding': typeof AuthOnboardingRouteRoute
  '/signin': typeof AuthSigninRouteRoute
  '/signup': typeof AuthSignupRouteRoute
  '/friend': typeof ProtectedLayoutFriendRouteRoute
  '/notification': typeof ProtectedLayoutNotificationRouteRoute
  '/call/$id': typeof ProtectedCallIdRouteRoute
  '/': typeof ProtectedLayoutIndexRoute
  '/chat/$id': typeof ProtectedLayoutChatIdRouteRoute
}

export interface FileRoutesById {
  __root__: typeof rootRoute
  '/_protected': typeof ProtectedRouteRouteWithChildren
  '/_auth/onboarding': typeof AuthOnboardingRouteRoute
  '/_auth/signin': typeof AuthSigninRouteRoute
  '/_auth/signup': typeof AuthSignupRouteRoute
  '/_protected/_layout': typeof ProtectedLayoutRouteRouteWithChildren
  '/_protected/_layout/friend': typeof ProtectedLayoutFriendRouteRoute
  '/_protected/_layout/notification': typeof ProtectedLayoutNotificationRouteRoute
  '/_protected/call/$id': typeof ProtectedCallIdRouteRoute
  '/_protected/_layout/': typeof ProtectedLayoutIndexRoute
  '/_protected/_layout/chat/$id': typeof ProtectedLayoutChatIdRouteRoute
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths:
    | ''
    | '/onboarding'
    | '/signin'
    | '/signup'
    | '/friend'
    | '/notification'
    | '/call/$id'
    | '/'
    | '/chat/$id'
  fileRoutesByTo: FileRoutesByTo
  to:
    | ''
    | '/onboarding'
    | '/signin'
    | '/signup'
    | '/friend'
    | '/notification'
    | '/call/$id'
    | '/'
    | '/chat/$id'
  id:
    | '__root__'
    | '/_protected'
    | '/_auth/onboarding'
    | '/_auth/signin'
    | '/_auth/signup'
    | '/_protected/_layout'
    | '/_protected/_layout/friend'
    | '/_protected/_layout/notification'
    | '/_protected/call/$id'
    | '/_protected/_layout/'
    | '/_protected/_layout/chat/$id'
  fileRoutesById: FileRoutesById
}

export interface RootRouteChildren {
  ProtectedRouteRoute: typeof ProtectedRouteRouteWithChildren
  AuthOnboardingRouteRoute: typeof AuthOnboardingRouteRoute
  AuthSigninRouteRoute: typeof AuthSigninRouteRoute
  AuthSignupRouteRoute: typeof AuthSignupRouteRoute
}

const rootRouteChildren: RootRouteChildren = {
  ProtectedRouteRoute: ProtectedRouteRouteWithChildren,
  AuthOnboardingRouteRoute: AuthOnboardingRouteRoute,
  AuthSigninRouteRoute: AuthSigninRouteRoute,
  AuthSignupRouteRoute: AuthSignupRouteRoute,
}

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>()

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/_protected",
        "/_auth/onboarding",
        "/_auth/signin",
        "/_auth/signup"
      ]
    },
    "/_protected": {
      "filePath": "_protected/route.tsx",
      "children": [
        "/_protected/_layout",
        "/_protected/call/$id"
      ]
    },
    "/_auth/onboarding": {
      "filePath": "_auth/onboarding/route.tsx"
    },
    "/_auth/signin": {
      "filePath": "_auth/signin/route.tsx"
    },
    "/_auth/signup": {
      "filePath": "_auth/signup/route.tsx"
    },
    "/_protected/_layout": {
      "filePath": "_protected/_layout/route.tsx",
      "parent": "/_protected",
      "children": [
        "/_protected/_layout/friend",
        "/_protected/_layout/notification",
        "/_protected/_layout/",
        "/_protected/_layout/chat/$id"
      ]
    },
    "/_protected/_layout/friend": {
      "filePath": "_protected/_layout/friend/route.tsx",
      "parent": "/_protected/_layout"
    },
    "/_protected/_layout/notification": {
      "filePath": "_protected/_layout/notification/route.tsx",
      "parent": "/_protected/_layout"
    },
    "/_protected/call/$id": {
      "filePath": "_protected/call/$id/route.tsx",
      "parent": "/_protected"
    },
    "/_protected/_layout/": {
      "filePath": "_protected/_layout/index.tsx",
      "parent": "/_protected/_layout"
    },
    "/_protected/_layout/chat/$id": {
      "filePath": "_protected/_layout/chat/$id/route.tsx",
      "parent": "/_protected/_layout"
    }
  }
}
ROUTE_MANIFEST_END */
