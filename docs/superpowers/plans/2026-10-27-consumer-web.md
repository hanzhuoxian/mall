# Consumer Web Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an independent responsive consumer storefront in `web/mall-web` using the existing marketplace API, adding the missing public product detail endpoint.

**Architecture:** Add one consumer-facing Umi Max application, separate from `mall-admin`, with its own API client, auth flow, routes, and storefront shell. Extend the marketplace service/store/controller with an ID-based catalog lookup that enforces the same publication and shop visibility rules as catalog listing.

**Tech Stack:** Go, Gin, GORM, React, TypeScript, Umi Max, Ant Design, Vitest, Testing Library.

**Spec:** `docs/superpowers/specs/2026-10-27-consumer-web-design.md`

## Global Constraints

- Consumers only; no merchant workbench or platform administration.
- All prices, stock, order state, and checkout totals are server-authoritative.
- No real payment integration; only show “开发沙箱付款” when `sandboxPayment` is true.
- Keep the new app under `web/mall-web`; do not import source files from `web/mall-admin`.
- Public catalog detail must hide pending/offline products and products whose shop is not approved.
- A checkout retry reuses its `requestKey` only for the same user, address ID, and cart ID set.

## Review Focus

- A deep link to a hidden or unavailable product must not disclose it; Task 1 tests pending product, offline product, and unapproved shop as 404.
- An expired access token with an invalid refresh token must clear session and preserve the intended login return path; Task 2 tests refresh success and refresh failure.
- Retrying checkout after a network timeout must reuse the same key and exact payload; Task 4 tests unchanged retry and changed address/cart ID set.
- A disabled sandbox must never expose a payment action; Task 5 tests sandbox off and on, including exact “开发沙箱付款” copy.
- Registration does not create an authenticated session; Task 2 tests success followed by explicit login and return navigation.

---

### Task 1: Public product detail API

**Files:**
- Modify: `internal/marketplace/store/store.go`
- Modify: `internal/marketplace/store/mysql/catalog.go`
- Modify: `internal/marketplace/service/query.go`
- Modify: `internal/marketplace/controller/product.go`
- Modify: `internal/marketplace/controller/routes.go`
- Test: `internal/marketplace/integration_test.go`

**Interfaces:**
- Produces: `Store.CatalogProduct(context.Context, string) (model.Product, error)` and `(*Service).CatalogProduct(context.Context, string) (model.Product, error)`.
- HTTP: `GET /v1/market/catalog/:id` returns the existing response envelope with one public `Product`; inaccessible or missing products return 404.

- [ ] **Step 1: Add a failing HTTP integration test** for an approved product (200 and product JSON), pending/offline product (404), and product at an unapproved shop (404).
- [ ] **Step 2: Run the focused test** with `go test ./internal/marketplace -run TestCatalogByIDVisibility -count=1` and confirm it fails because the route/lookup is missing.
- [ ] **Step 3: Implement `CatalogProduct`** through the Store interface and MySQL store. Use the catalog predicate for approved products and approved shops; translate no row to `store.ErrNotFound`.
- [ ] **Step 4: Register `GET /catalog/:id`** and return the service result through the existing `respond` helper.
- [ ] **Step 5: Re-run the focused test** and `go test ./internal/marketplace/controller ./internal/marketplace/store/mysql`; confirm visible product succeeds and all hidden cases are 404.

### Task 2: Independent app shell and account session

**Files:**
- Create: `web/mall-web/package.json`, `package-lock.json`, `tsconfig.json`, `.editorconfig`, `biome.json`
- Create: `web/mall-web/config/config.ts`, `proxy.ts`, `routes.ts`
- Create: `web/mall-web/src/app.tsx`, `global.less`, `request.ts`
- Create: `web/mall-web/src/layouts/StoreLayout.tsx`
- Create: `web/mall-web/src/pages/user/login/index.tsx`, `register/index.tsx`
- Create: `web/mall-web/src/pages/user/auth.test.tsx`
- Create: `web/mall-web/src/pages/user/service.ts`

**Interfaces:**
- Produces: auth service `login(credentials)`, `register(input)`, `logout()`, `refreshSession()`, and request helper `api<T>(path, options)`.
- `api<T>` consumes raw `{ code, message, data }`, adds Bearer access token, and throws a typed API error for non-2xx or nonzero code.

- [ ] **Step 1: Scaffold the independent Umi Max app and test harness** with local package metadata, routes/config, and dev/test/build scripts aligned to the versions in `web/mall-admin`; do not add auth behavior yet.
- [ ] **Step 2: Add failing tests** in `src/pages/user/auth.test.tsx` for raw envelope parsing, Bearer header, successful refresh/retry, failed refresh clearing both tokens, login return path, and registration requiring explicit login.
- [ ] **Step 3: Run** `cd web/mall-web && npm test -- --run src/pages/user/auth.test.tsx`; confirm failures before implementing auth behavior.
- [ ] **Step 4: Implement the request/auth services** using `/login`, `/register`, `/captcha`, `/refresh`, and `/logout`; store returned access/refresh tokens in localStorage and retry a failed request at most once after refresh.
- [ ] **Step 5: Implement responsive login/register pages**; registration uses captcha and success returns to login, while login restores a safe same-origin `redirect` path.
- [ ] **Step 6: Re-run auth tests and app typecheck** with `npm test -- --run src/pages/user/auth.test.tsx` and `npm run tsc`; confirm pass.

### Task 3: Storefront catalog and product detail

**Files:**
- Create: `web/mall-web/src/pages/home/index.tsx`, `home.test.tsx`
- Create: `web/mall-web/src/pages/product/detail.tsx`, `detail.test.tsx`
- Create: `web/mall-web/src/services/marketplace.ts`
- Modify: `web/mall-web/config/routes.ts`
- Modify: `web/mall-web/src/layouts/StoreLayout.tsx`, `global.less`

**Interfaces:**
- Consumes Task 1 `GET /catalog/:id` and Task 2 `api<T>`.
- Produces marketplace client methods `catalog(params)`, `categories(params)`, `product(id)`, and `putCart(productId, quantity)`.

- [ ] **Step 1: Add failing page tests** for public categories/catalog rendering, name/category filters, pagination, direct product detail load, unavailable-product state, and add-to-cart login redirect.
- [ ] **Step 2: Run** `cd web/mall-web && npm test -- --run src/pages/home/home.test.tsx src/pages/product/detail.test.tsx`; confirm expected failures.
- [ ] **Step 3: Implement typed marketplace calls** with `name`, `categoryId`, `current`, and `pageSize` query parameters.
- [ ] **Step 4: Implement the responsive store shell, catalog, and product detail pages** using server data only; use browser-safe image placeholders when product image URL is absent.
- [ ] **Step 5: Re-run page tests and TypeScript check**; confirm filtering, direct links, and login redirect work.

### Task 4: Cart, addresses, and idempotent checkout

**Files:**
- Create: `web/mall-web/src/pages/cart/index.tsx`, `cart.test.tsx`
- Create: `web/mall-web/src/pages/checkout/index.tsx`, `checkout.test.tsx`
- Create: `web/mall-web/src/pages/account/addresses.tsx`, `addresses.test.tsx`
- Modify: `web/mall-web/src/services/marketplace.ts`

**Interfaces:**
- Consumes Task 2 authenticated request/session layer and Task 3 marketplace types.
- Produces calls for `GET/PUT/DELETE /cart`, `GET/POST/DELETE /addresses`, and `POST /checkout` with `{ addressId, cartIds, requestKey }`.

- [ ] **Step 1: Add failing tests** for quantity update/delete, selection, address create/delete, server total display, checkout success, validation failure, and retry-key rules.
- [ ] **Step 2: Run focused cart/checkout/address tests** and confirm failures.
- [ ] **Step 3: Implement cart and address service methods** using current-user-owned endpoints and existing API shapes.
- [ ] **Step 4: Implement the cart, address, and checkout pages**; persist an in-flight checkout signature/key so an unchanged retry reuses the key, and generate a new key if address or selected cart IDs change.
- [ ] **Step 5: Re-run focused tests and TypeScript check**; confirm errors preserve user selections and checkout uses server-returned orders/totals.

### Task 5: Consumer orders, sandbox copy, and release checks

**Files:**
- Create: `web/mall-web/src/pages/orders/index.tsx`, `orders.test.tsx`
- Modify: `web/mall-web/src/services/marketplace.ts`
- Modify: `web/mall-web/config/proxy.ts`, `routes.ts`
- Modify: `docs/接口与验收.md`

**Interfaces:**
- Consumes order list (including `items`), authenticated context (`sandboxPayment`), and Task 2 request/session layer.
- Produces state-appropriate consumer actions: cancel, “开发沙箱付款”, receive, and request refund.

- [ ] **Step 1: Add failing tests** for buyer order filtering/details and permitted actions per state; assert payment action absent with sandbox off and exact “开发沙箱付款” label with sandbox on.
- [ ] **Step 2: Run** `cd web/mall-web && npm test -- --run src/pages/orders/orders.test.tsx`; confirm failures.
- [ ] **Step 3: Implement orders list/detail and actions** from `GET /orders?scope=buyer` and `POST /orders/:id/actions`; keep the payment disclosure visible before and after the sandbox action.
- [ ] **Step 4: Configure dev/test proxies** for `/v1/`, `/login`, `/register`, `/logout`, `/refresh`, and `/captcha`; document start/build commands for `web/mall-web`.
- [ ] **Step 5: Run full verification**: `go test ./internal/marketplace/...`, then in `web/mall-web` run `npm test`, `npm run tsc`, and `npm run build`; confirm all pass.
