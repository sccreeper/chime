import { get, writable } from "svelte/store";

const notification_max = Math.pow(2, 16);


export const notifications = writable([{id: null, component: null, props: null}]);

// Notification utility functions

export function removeNotification(id) {
    
    let index = 0;

    for (let i = 0; i < get(notifications).length; i++) {
        
        if (get(notifications)[i].id == id) {
            index = i
            break
        }

    }

    let n = get(notifications)
    n.splice(index, 1)
    notifications.set(n)

}

export function createNotification(component, props) {


    let id = notificationID()
    let p = props
    p.id = id

    let n = get(notifications)

    n.push({id: id, component: component, props: p})
    notifications.set(n)

}

export function notificationID() {
    
    let id;
    id = Math.floor(Math.random() * notification_max)

    if (get(notifications).length == 0) {
        return id
    }

    for (let i = 0; i < get(notifications).length; i++) {

        if (get(notifications)[i].id == id) {
            id = Math.floor(Math.random() * notification_max)
            continue
        } else {
            return id
        }
        
    }

}