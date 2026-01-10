export interface UserInfo {
  SID: string;
  UserName: string;
  IsEnter: boolean;
}

export interface CardMessage {
  IsCard: boolean;
  IsNew: boolean;
  CardID?: string;
  SID?: string;
}

export interface RegistResponse {
  Success: boolean;
}

export interface QuestionData {
  Use: string[];
  message: string;
}

export interface ErrorInfo {
  message: string;
}