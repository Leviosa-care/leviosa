import { persisted } from 'svelte-persisted-store'

import type { NavState, ConsultationState, EventState, MessageState, ReservationState, ServiceState } from '$lib/types'
import { NAV_STATES, CONSULTATION_STATES, EVENT_STATES, MESSAGE_STATES, RESERVATION_STATES, SERVICE_STATES } from '$lib/types'

const navigationStateInit = $state<NavState>(NAV_STATES.Accueil)
const consultationStateInit = $state<ConsultationState>(CONSULTATION_STATES.ConsultationsAVenir)
const eventStateInit = $state<EventState>(EVENT_STATES.EvenementsAVenir)
const messageStateInit = $state<MessageState>(MESSAGE_STATES.NotesDeSeances)
const reservationStateInit = $state<ReservationState>(RESERVATION_STATES.Consultations)
const serviceStateInit = $state<ServiceState>(SERVICE_STATES.APropos)

export const navigationState = persisted('navigationState', navigationStateInit)
export const consultationState = persisted('consultationState', consultationStateInit)
export const eventState = persisted('eventState', eventStateInit)
export const messageState = persisted('messageState', messageStateInit)
export const reservationState = persisted('reservationState', reservationStateInit)
export const serviceState = persisted('serviceState', serviceStateInit)

