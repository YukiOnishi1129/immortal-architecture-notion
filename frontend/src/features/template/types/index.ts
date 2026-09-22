export interface TemplateField {
  id: string;
  label: string;
  order: number;
  isRequired: boolean;
}

export interface TemplateOwner {
  id: string;
  firstName: string;
  lastName: string;
  thumbnail: string | null;
}

export interface Template {
  id: string;
  name: string;
  ownerId?: string;
  owner?: TemplateOwner;
  fields: TemplateField[];
  isUsed?: boolean;
  createdAt?: string;
  updatedAt: string;
  /** Notionの置き場所となる親ページのURL。未連携なら null。 */
  notionParentPageUrl?: string | null;
}

export interface TemplateFilters {
  q?: string;
  page?: number;
  ownerId?: string;
  onlyMyTemplates?: boolean;
}
