// node modules
import { action, makeObservable, observable } from 'mobx';
import * as R from 'ramda';

export default class ContactsBlockStore {
  // properties
  @observable name: string = '';
  @observable phone: string = '';
  @observable email: string = '';
  @observable description: string = '';

  @observable errors: any = {};

  constructor() {
    makeObservable(this);
  }
  
  // actions
  @action setName = (data: string) => {
    this.name = data;
  };

  @action setPhone = (data: string) => {
    this.phone = data;
  };

  @action setEmail = (data: string) => {
    this.email = data;
  };

  @action setDescription = (data: string) => {
    this.description = data;
  };

  @action setErrors = (key: string, errorText: string) => {
    this.errors = {
      ...this.errors,
      [key]: errorText,
    };
  };

  @action clearErrorByKey = (key: string) => {
    this.errors = R.dissoc(key, this.errors);
  };
}
