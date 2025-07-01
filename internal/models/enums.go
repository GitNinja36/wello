package models

type Role string

const (
	PATIENT Role = "PATIENT"
	DOCTOR  Role = "DOCTOR"
	ADMIN   Role = "ADMIN"
)

type AppointmentMode string

const (
	APPT_MODE_ONLINE  AppointmentMode = "ONLINE"
	APPT_MODE_OFFLINE AppointmentMode = "OFFLINE"
)

type AppointmentStatus string

const (
	PENDING               AppointmentStatus = "PENDING"
	ACCEPTED              AppointmentStatus = "ACCEPTED"
	REJECTED              AppointmentStatus = "REJECTED"
	COMPLETED             AppointmentStatus = "COMPLETED"
	RESCHEDULE_REQUESTED  AppointmentStatus = "RESCHEDULE_REQUESTED"
	RESCHEDULED           AppointmentStatus = "RESCHEDULED"
	RESCHEDULE_REJECTED   AppointmentStatus = "RESCHEDULE_REJECTED"
	RESCHEDULED_CONFIRMED AppointmentStatus = "RESCHEDULED_CONFIRMED"
	CANCELLED_BY_PATIENT  AppointmentStatus = "CANCELLED_BY_PATIENT"
)

type TestType string

const (
	BLOOD TestType = "BLOOD"
	URINE TestType = "URINE"
	XRAY  TestType = "XRAY"
	CT    TestType = "CT"
	OTHER TestType = "OTHER"
)

type TestStatus string

const (
	TEST_PENDING TestStatus = "PENDING"
	SCHEDULED    TestStatus = "SCHEDULED"
	TEST_DONE    TestStatus = "COMPLETED"
	REPORTED     TestStatus = "REPORTED"
)

type OrderStatus string

const (
	ORDER_PENDING   OrderStatus = "PENDING"
	PROCESSING      OrderStatus = "PROCESSING"
	SHIPPED         OrderStatus = "SHIPPED"
	DELIVERED       OrderStatus = "DELIVERED"
	ORDER_CANCELLED OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PAYMENT_ONLINE PaymentMethod = "ONLINE"
	PAYMENT_COD    PaymentMethod = "COD"
)
