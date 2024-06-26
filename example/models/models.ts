// Generated by sqlboiler-erg: DO NOT EDIT.
export enum Role {
  SUPER_ADMIN = "SUPER_ADMIN",
  ADMIN = "ADMIN",
  EMPLOYEE = "EMPLOYEE",
}
export interface AccountRelations {
  accountSessions: AccountSession[];
  comments: Comment[];
  events: Event[];
  organizationAccounts: OrganizationAccount[];
}

export interface Account extends AccountRelations {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;

  customFields?: Record<string, any>
}
export interface AccountSessionRelations {
  account?: Account;
}

export interface AccountSession extends AccountSessionRelations {
  id: string;
  accountId: string;
  ipAddress: string;
  userAgent: string;
  expiresAt: Date;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;

  customFields?: Record<string, any>
}
export interface CommentRelations {
  account?: Account;
  event?: Event;
  organization?: Organization;
}

export interface Comment extends CommentRelations {
  id: string;
  organizationId: string;
  eventId: string;
  accountId: string;
  comment: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;

  customFields?: Record<string, any>
}
export interface EventRelations {
  account?: Account;
  organization?: Organization;
  comments: Comment[];
}

export interface Event extends EventRelations {
  id: string;
  organizationId: string;
  accountId: string;
  name: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;

  customFields?: Record<string, any>
}
export interface OrganizationRelations {
  comments: Comment[];
  events: Event[];
  organizationAccounts: OrganizationAccount[];
}

export interface Organization extends OrganizationRelations {
  id: string;
  name: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;

  customFields?: Record<string, any>
}
export interface OrganizationAccountRelations {
  account?: Account;
  organization?: Organization;
}

export interface OrganizationAccount extends OrganizationAccountRelations {
  id: string;
  organizationId: string;
  accountId: string;
  role: Role;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;

  customFields?: Record<string, any>
}