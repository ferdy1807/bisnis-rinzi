# Bisnis-Rinzi

<p align="center">
  <img src="https://img.shields.io/badge/Backend-Go%201.24-blue?logo=go">
  <img src="https://img.shields.io/badge/Frontend-Vue%203-green?logo=vuedotjs">
  <img src="https://img.shields.io/badge/Database-PostgreSQL-blue?logo=postgresql">
  <img src="https://img.shields.io/badge/Container-Docker-2496ED?logo=docker">
  <img src="https://img.shields.io/badge/Architecture-Microservices-orange">
  <img src="https://img.shields.io/badge/PWA-Offline%20First-purple">
  <img src="https://img.shields.io/badge/Cache-Redis-red?logo=redis">
  <img src="https://img.shields.io/badge/Storage-MinIO-orange?logo=minio">
</p>

## Overview

**Bisnis-Rinzi** is a modern business management platform built using a **Microservice Architecture** and **Clean Architecture** principles.

The system is designed to support:

- Inventory Management
- Point of Sale (POS)
- Cash Management
- Rental Management
- Financial Management
- Authentication & Single Sign-On (SSO)

---

# Technology Stack

## Backend

- Go (Golang)
- Clean Architecture
- REST API
- API Gateway
- PostgreSQL
- Redis Streams
- Outbox Pattern
- MinIO Object Storage

## Frontend

- Vue 3
- Vite
- Pinia
- Vue Router
- PrimeVue
- Progressive Web App (PWA)
- IndexedDB Offline Storage

## Infrastructure

- Docker
- Docker Compose
- PostgreSQL
- Redis
- MinIO

---

# System Architecture

```text
                    +------------------+
                    |   API Gateway    |
                    +---------+--------+
                              |
      ---------------------------------------------------
      |          |           |          |          |
      v          v           v          v          v

+-----------+ +-----------+ +--------+ +---------+ +-----------+
|   Auth    | | Inventory | |  POS   | | Rental  | | Finance   |
|  Service  | |  Service  | |Service | | Service | | Service   |
+-----------+ +-----------+ +--------+ +---------+ +-----------+

      |           |           |           |           |
      +-----------+-----------+-----------+-----------+
                              |
                       Redis Streams
                         Event Bus

      -------------------------------------------------
      |               Infrastructure                  |
      -------------------------------------------------

      PostgreSQL (Database per Service)
      Redis Streams
      MinIO Object Storage
```

---

# Portal TOKO

![alt text](foto/portal-toko.png)

# =====================================

# Admin Dashboard

![alt text](foto/admin-dashboard.png)

# =====================================

# Portal Sewa Hantaran

![alt text](foto/portal-sewa.png)

# =====================================

- Internal Gateway

## Microservices

- Auth Service
- Inventory Service
- POS Service
- Rental Service
- Finance Service
- Cash Service

## Auth Service

Responsible for:

- Authentication
- Authorization
- JWT Management
- Single Sign-On (SSO)
- User Management
- Role Management

## Event Driven

- Go outbox_event
