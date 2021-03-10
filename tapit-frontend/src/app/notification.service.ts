import { Injectable } from '@angular/core';

export class Notification {
  id: number;
  resultType: string; // enum success or failure or info
  text: string;
}

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  notifications: Notification[] = [];
  currentCount = 0;

  addNotification(resultType, text) {
    const newNotification = new Notification();
    newNotification.id = this.currentCount;
    this.currentCount++;
    newNotification.resultType = resultType;
    newNotification.text = text;

    this.notifications.push(newNotification);
    setTimeout(() => this.closeNotification(newNotification), 3000);
  }

  closeNotification(notify: Notification) {
    for (let i = 0; i < this.notifications.length; i++) {
      if (this.notifications[i].id === notify.id) {
        this.notifications.splice(i, 1);
        break;
      }
    }
  }

  constructor() {
  }
}
