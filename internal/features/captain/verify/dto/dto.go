package dto

import activity_verfication "encore.app/background/temporal/captain/verification/activity"

type UploadDocumentRequest struct {
	FileData    string `json:"file_data" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}

type UploadDocumentResponse struct {
	Message string `json:"message"`
}

type VerificationStatusResponse = activity_verfication.BackgroundVerificationWorkflowInput
