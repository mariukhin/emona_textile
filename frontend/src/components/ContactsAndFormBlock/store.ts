// node modules
import { action, makeObservable, observable } from 'mobx';
// mocks

export default class ContactsBlockStore {
  // properties
  @observable name: string | null = null;
  @observable phone: string | null = null;
  @observable email: string | null = null;
  @observable description: string | null = null;

  @observable formIsValid: boolean = true;
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

  @action setFormIsValid = (data: boolean) => {
    this.formIsValid = data;
  };

  @action setErrors = (key: string, errorText: string) => {
    this.errors = {
      ...this.errors,
      [key]: errorText,
    };
  };
}
