import { templatesApiClient } from "@/external/client/api/config";
import type { TemplatesApi } from "@/external/client/api/generated/apis/TemplatesApi";
import type { ModelsTemplateResponse } from "@/external/client/api/generated/models/ModelsTemplateResponse";
import {
  type CreateTemplateRequest,
  type TemplateResponse,
  TemplateResponseSchema,
  type UpdateTemplateRequest,
} from "@/external/dto/template.dto";
import { isNotFoundError } from "../http-error";

function toTemplateResponse(model: ModelsTemplateResponse): TemplateResponse {
  return TemplateResponseSchema.parse({
    id: model.id,
    name: model.name,
    ownerId: model.ownerId,
    owner: {
      id: model.owner.id,
      firstName: model.owner.firstName,
      lastName: model.owner.lastName,
      thumbnail: model.owner.thumbnail ?? null,
    },
    fields: model.fields.map((field) => ({
      id: field.id,
      label: field.label,
      order: field.order,
      isRequired: field.isRequired,
    })),
    updatedAt: model.updatedAt.toISOString(),
    isUsed: model.isUsed,
  });
}

export class TemplateService {
  constructor(private readonly api: TemplatesApi) {}

  async getTemplateById(id: string): Promise<TemplateResponse | null> {
    try {
      const template = await this.api.templatesGetTemplateById({
        templateId: id,
      });
      return toTemplateResponse(template);
    } catch (error) {
      if (isNotFoundError(error)) {
        return null;
      }
      throw error;
    }
  }

  async getTemplates(filters?: {
    ownerId?: string;
    q?: string;
  }): Promise<TemplateResponse[]> {
    const templates = await this.api.templatesListTemplates({
      ownerId: filters?.ownerId,
      q: filters?.q,
    });
    return templates.map((template) => toTemplateResponse(template));
  }

  async createTemplate(
    ownerId: string,
    input: CreateTemplateRequest,
  ): Promise<TemplateResponse> {
    const template = await this.api.templatesCreateTemplate({
      modelsCreateTemplateRequest: {
        name: input.name,
        ownerId,
        fields: input.fields.map((field) => ({
          label: field.label,
          order: field.order,
          isRequired: field.isRequired,
        })),
      },
    });
    return toTemplateResponse(template);
  }

  async updateTemplate(
    id: string,
    ownerId: string,
    input: UpdateTemplateRequest,
  ): Promise<TemplateResponse> {
    const template = await this.api.templatesUpdateTemplate({
      templateId: id,
      ownerId,
      modelsUpdateTemplateRequest: {
        id,
        name: input.name,
        fields: input.fields.map((field) => ({
          id: field.id,
          label: field.label,
          order: field.order,
          isRequired: field.isRequired,
        })),
      },
    });
    return toTemplateResponse(template);
  }

  async deleteTemplate(id: string, ownerId: string): Promise<void> {
    await this.api.templatesDeleteTemplate({ templateId: id, ownerId });
  }
}

export const templateService = new TemplateService(templatesApiClient);
