/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package handler

type Handler interface {
	Exec() error
	UpdateProgress(progress float32, status string) error
}
