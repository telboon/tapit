import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MainComponent } from './main/main.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { CampaignComponent } from './campaign/campaign.component';
import { CampaignNewComponent } from './campaign-new/campaign-new.component';
import { CampaignViewComponent } from './campaign-view/campaign-view.component';
import { PhonebookComponent } from './phonebook/phonebook.component';
import { PhonebookNewComponent } from './phonebook-new/phonebook-new.component';
import { TextTemplateComponent } from './text-template/text-template.component';
import { TextTemplateNewComponent } from './text-template-new/text-template-new.component';
import { ProviderComponent } from './provider/provider.component';
import { ProfileComponent } from './profile/profile.component';
import { WebTemplateComponent } from './web-template/web-template.component';
import { WebTemplateNewComponent } from './web-template-new/web-template-new.component';
import { GlobalSettingsComponent } from './global-settings/global-settings.component';

const routes: Routes = [
  { path: '', component: MainComponent },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'profile', component: ProfileComponent },
  { path: 'campaign', component: CampaignComponent },
  { path: 'campaign/new', component: CampaignNewComponent },
  { path: 'campaign/:id/view', component: CampaignViewComponent },
  { path: 'phonebook', component: PhonebookComponent },
  { path: 'phonebook/new', component: PhonebookNewComponent },
  { path: 'phonebook/:id/edit', component: PhonebookNewComponent },
  { path: 'text-template', component: TextTemplateComponent },
  { path: 'text-template/new', component: TextTemplateNewComponent },
  { path: 'text-template/:id/edit', component: TextTemplateNewComponent },
  { path: 'provider', component: ProviderComponent },
  { path: 'web-template', component: WebTemplateComponent },
  { path: 'web-template/new', component: WebTemplateNewComponent },
  { path: 'web-template/:id/edit', component: WebTemplateNewComponent },
  { path: 'global-settings', component: GlobalSettingsComponent },
  ];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
